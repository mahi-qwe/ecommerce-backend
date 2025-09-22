import React from "react";
import { Link, useLocation } from "react-router-dom";
import { LayoutDashboard, Users, Package, ShoppingCart } from "lucide-react";

export default function Sidebar() {
  const location = useLocation();

  // Function to determine if a link is active
  const isActive = (path) => location.pathname === path;

  // Navigation items with icons
  const navItems = [
    {
      path: "/dashboard",
      label: "Dashboard",
      icon: LayoutDashboard,
      description: "Overview & Analytics",
    },
    {
      path: "/users",
      label: "Users",
      icon: Users,
      description: "Manage Users",
    },
    {
      path: "/products",
      label: "Products",
      icon: Package,
      description: "Product Catalog",
    },
    {
      path: "/orders",
      label: "Orders",
      icon: ShoppingCart,
      description: "Order Management",
    },
  ];

  return (
    <aside className="w-64 min-h-screen bg-gradient-to-b from-slate-800 to-slate-900 border-r border-slate-700 shadow-xl">
      <div className="p-6">
        {/* Sidebar Header */}
        <div className="mb-8">
          <h2 className="text-xl font-bold text-white mb-2">Admin Panel</h2>
          <div className="h-1 w-12 bg-blue-500 rounded-full"></div>
        </div>

        {/* Navigation Menu */}
        <nav className="space-y-2">
          {navItems.map((item) => {
            const Icon = item.icon;
            const active = isActive(item.path);

            return (
              <Link
                key={item.path}
                to={item.path}
                className={`
                  group flex items-center px-4 py-3 rounded-xl transition-all duration-200
                  ${
                    active
                      ? "bg-blue-500 text-white shadow-lg transform scale-105"
                      : "text-gray-300 hover:bg-slate-700 hover:text-white hover:transform hover:scale-102"
                  }
                `}
              >
                {/* Icon */}
                <Icon
                  className={`
                    w-5 h-5 mr-3 transition-colors duration-200
                    ${
                      active
                        ? "text-white"
                        : "text-gray-400 group-hover:text-white"
                    }
                  `}
                />

                {/* Label and Description */}
                <div className="flex-1">
                  <div className="font-medium text-sm">{item.label}</div>
                  <div
                    className={`
                    text-xs transition-colors duration-200
                    ${
                      active
                        ? "text-blue-100"
                        : "text-gray-500 group-hover:text-gray-300"
                    }
                  `}
                  >
                    {item.description}
                  </div>
                </div>

                {/* Active indicator */}
                {active && (
                  <div className="w-2 h-2 bg-white rounded-full"></div>
                )}
              </Link>
            );
          })}
        </nav>

        {/* Bottom decoration
        <div className="mt-8 pt-8 border-t border-slate-700">
          <div className="text-center">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full mx-auto mb-2"></div>
            <p className="text-xs text-gray-400">Admin Tools</p>
          </div>
        </div> */}
      </div>
    </aside>
  );
}
