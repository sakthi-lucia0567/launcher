// src/components/AppList.jsx
import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

import { isMobile, isWindows } from "react-device-detect";

const AppList = () => {
  const [apps, setApps] = useState([]);
  const [loading, setLoading] = useState(false);
  const [currentApp, setCurrentApp] = useState(null);

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

  const launchApp = (app) => {
    setLoading(true);
    axios
      .post(import.meta.env.VITE_SERVER_URL + "/v1/launcher", {
        name: app.name,
        path: app.path,
      })
      .then(() => {
        setLoading(false);
        setCurrentApp(app);
      })
      .catch((error) => {
        setLoading(false);
        console.error("Error launching app:", error);
      });
  };

  const quitApp = () => {
    setLoading(true);
    axios
      .post(import.meta.env.VITE_SERVER_URL + "/v1/quit", {
        name: currentApp.name,
        application: currentApp.application,
      })
      .then(() => {
        setLoading(false);
        setCurrentApp(null);
        fetchApps();
      })
      .catch((error) => {
        setLoading(false);
        console.error("Error quitting app:", error);
      });
  };

  if (loading) {
    return <div className="text-center mt-8">Loading...</div>;
  }

  if (currentApp) {
    return (
      <div className="container mx-auto px-4 py-8 text-center">
        <h1 className="text-3xl font-bold mb-6">
          Application Opened: {currentApp.name}
        </h1>
        {isWindows && ( // Render settings button only if showSettings is true
          <a
            onClick={() => navigate("/settings")}
            className="absolute bottom-4 left-4 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded"
          >
            Settings
          </a>
        )}
        <button
          onClick={quitApp}
          className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded"
        >
          Home
        </button>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6 text-center">Applications</h1>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        {isWindows && ( // Render settings button only if showSettings is true
          <button
            onClick={() => navigate("/settings")}
            className="absolute bottom-4 left-4 bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded"
          >
            Settings
          </button>
        )}
        {apps.map((app) => (
          <div
            key={app.id}
            onClick={() => launchApp(app)}
            className="bg-white rounded-lg shadow-md p-4 flex flex-col items-center transition-transform hover:scale-105"
          >
            <img
              src={import.meta.env.VITE_IMAGE_URL + app.icon}
              alt={app.name}
              className="w-16 h-16 mb-2"
            />
            <h2 className="text-center font-semibold">{app.name}</h2>
          </div>
        ))}
      </div>
    </div>
  );
};

export default AppList;
