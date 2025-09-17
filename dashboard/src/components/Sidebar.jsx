import React from "react";
import { Link } from "react-router-dom";

const Sidebar = () => {
  return (
    <div className="h-screen w-60 bg-gray-900 text-white fixed left-0 top-0 p-5">
      <h2 className="text-2xl font-bold mb-10">Admin Panel</h2>
      <nav className="flex flex-col gap-4">
        <Link to="/dashboard" className="hover:bg-gray-700 p-2 rounded">
          Dashboard
        </Link>
        <Link to="/users" className="hover:bg-gray-700 p-2 rounded">
          Users
        </Link>
        <Link to="/products" className="hover:bg-gray-700 p-2 rounded">
          Products
        </Link>
        <Link to="/orders" className="hover:bg-gray-700 p-2 rounded">
          Orders
        </Link>
      </nav>
    </div>
  );
};

export default Sidebar;
