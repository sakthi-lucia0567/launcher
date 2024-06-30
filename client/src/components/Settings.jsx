// src/components/Settings.jsx
import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const Settings = () => {
  const [apps, setApps] = useState([]);
  const [selectedApp, setSelectedApp] = useState(null);
  const [appConfig, setAppConfig] = useState({});
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newApp, setNewApp] = useState({ name: "", path: "", icon: null });
  const navigate = useNavigate();

  useEffect(() => {
    fetchApps();
  }, []);

  const fetchApps = () => {
    axios
      .get(import.meta.env.VITE_SERVER_URL + "/v1/get_application")
      .then((response) => {
        setApps(response.data);
      });
  };

  const handleAppClick = (app) => {
    console.log("app", app);
    setSelectedApp(app);
    setAppConfig(app.path);
  };

  const handleInputChange = (event) => {
    setAppConfig(event.target.value);
  };

  const saveConfig = (id) => {
    console.log("id is clicked", id);
    axios
      .put(`${import.meta.env.VITE_SERVER_URL}/v1/update_application/${id}`, {
        path: appConfig,
      })
      .then(() => {
        fetchApps();
        setSelectedApp(null);
        setAppConfig({});
      })
      .catch((error) => {
        console.error("Error updating config:", error);
      });
  };

  const handleNewAppChange = (event) => {
    const { name, value, files } = event.target;
    setNewApp((prev) => ({
      ...prev,
      [name]: name === "icon" ? files[0] : value,
    }));
  };

  const handleSubmitNewApp = async (e) => {
    e.preventDefault();
    try {
      // Upload the file and get the file name in the response
      const formData = new FormData();
      formData.append("file", newApp.icon);
      const uploadResponse = await axios.post(
        import.meta.env.VITE_SERVER_URL + "/v1/upload",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      );
      console.log("boke 1", uploadResponse);

      const fileName = "/static/images/" + uploadResponse.data;

      const applicationData = {
        name: newApp.name,
        path: newApp.path,
        icon: fileName,
      };

      await axios.post(
        import.meta.env.VITE_SERVER_URL + "/v1/create_application",
        applicationData
      );

      setIsModalOpen(false);
      setNewApp({
        name: "",
        path: "",
        icon: null,
      });
      fetchApps();
    } catch (error) {
      console.error("Error creating application:", error);
    }
  };

  const removeApp = (id) => {
    // Logic to remove the selected app
    console.log("passed id", id);
    axios
      .delete(`${import.meta.env.VITE_SERVER_URL}/v1/delete_application/${id}`)
      .then(() => {
        fetchApps();
        setSelectedApp(null);
      })
      .catch((error) => {
        console.error("Error deleting app:", error);
      });
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6 text-center">Settings</h1>
      <button
        onClick={() => navigate("/")}
        className="mb-4 bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
      >
        Back to Home
      </button>
      <div className="divide-y divide-gray-300 border-2">
        {apps.map((app) => (
          <div
            key={app.id}
            onClick={() => handleAppClick(app)}
            className="py-4 px-2 flex justify-between items-center cursor-pointer hover:bg-gray-100 transition-colors"
          >
            <h2 className="text-lg font-semibold">{app.name}</h2>
          </div>
        ))}
      </div>
      <div className="flex justify-end items-end mt-4">
        <button
          onClick={() => setIsModalOpen(true)}
          className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
        >
          Add
        </button>
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
          <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
            <h3 className="text-lg font-bold mb-4">Add New Application</h3>
            <form onSubmit={handleSubmitNewApp}>
              <input
                type="text"
                name="name"
                value={newApp.name}
                onChange={handleNewAppChange}
                placeholder="Application Name"
                className="w-full p-2 mb-4 border border-gray-300 rounded"
              />
              <input
                type="text"
                name="path"
                value={newApp.path}
                onChange={handleNewAppChange}
                placeholder="Config Path"
                className="w-full p-2 mb-4 border border-gray-300 rounded"
              />
              <input
                type="file"
                name="icon"
                onChange={handleNewAppChange}
                className="w-full p-2 mb-4 border border-gray-300 rounded"
              />
              <div className="flex justify-end">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="mr-2 px-4 py-2 bg-gray-300 text-gray-800 rounded"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-500 text-white rounded"
                >
                  Submit
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
      <div className="mt-8">
        {selectedApp && (
          <>
            <h2 className="text-2xl font-bold mb-4">
              {selectedApp.name} Config
            </h2>
            <div className="mb-4">
              <label className="block text-gray-700">App Config Path</label>
              <input
                type="text"
                name="path"
                value={appConfig}
                onChange={handleInputChange}
                className="w-full p-2 border border-gray-300 rounded"
              />
            </div>
            <button
              onClick={() => saveConfig(selectedApp.id)}
              className="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded mr-4"
            >
              Save
            </button>

            <button
              onClick={() => removeApp(selectedApp.id)}
              className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded"
            >
              Remove
            </button>
          </>
        )}
      </div>
    </div>
  );
};

export default Settings;
