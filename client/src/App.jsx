import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AppList from "./components/AppList";
import Settings from "./components/Settings";

const App = () => {
  return (
    <div className="min-h-screen bg-gray-100">
      <Router>
        <Routes>
          <Route path="/" element={<AppList />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Router>
    </div>
  );
};

export default App;
