import React from "react";
import { Link, useLocation } from "react-router-dom";

export default function Sidebar() {
  const location = useLocation();

  const linkClasses = (path) =>
    `block px-4 py-2 rounded-lg hover:bg-gray-700 transition ${
      location.pathname === path ? "bg-gray-700 text-white" : "text-gray-300"
    }`;

  return (
    <aside className="w-64 min-h-screen bg-gray-800 p-4">
      <h2 className="text-xl font-bold text-white mb-6">Admin Panel</h2>
      <nav className="space-y-2">
        <Link to="/dashboard" className={linkClasses("/dashboard")}>
          Dashboard
        </Link>
        <Link to="/users" className={linkClasses("/users")}>
          Users
        </Link>
        <Link to="/products" className={linkClasses("/products")}>
          Products
        </Link>
        <Link to="/orders" className={linkClasses("/orders")}>
          Orders
        </Link>
      </nav>
    </aside>
  );
}
