import React from "react";
import { useNavigate } from "react-router-dom";
import { LogOut, User } from "lucide-react";

export default function Header() {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token"); // clear JWT
    navigate("/"); // go back to login
  };

  return (
    <header className="bg-gradient-to-r from-slate-800 to-slate-900 shadow-lg border-b border-slate-700">
      <div className="flex justify-between items-center px-6 py-4">
        {/* Left side - Logo and Title */}
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center">
            <User className="w-5 h-5 text-white" />
          </div>
          <h1 className="text-xl font-bold text-white tracking-wide">
            Admin Dashboard
          </h1>
        </div>

        {/* Right side - User info and logout */}
        <div className="flex items-center space-x-4">
          <div className="text-right hidden sm:block">
            <p className="text-sm text-gray-300">Welcome back</p>
            <p className="text-xs text-gray-400">Administrator</p>
          </div>

          <button
            onClick={handleLogout}
            className="flex items-center space-x-2 bg-red-500 hover:bg-red-600 px-4 py-2 rounded-lg transition-all duration-200 transform hover:scale-105 shadow-md"
          >
            <LogOut className="w-4 h-4 text-white" />
            <span className="text-white font-medium hidden sm:inline">
              Logout
            </span>
          </button>
        </div>
      </div>
    </header>
  );
}
