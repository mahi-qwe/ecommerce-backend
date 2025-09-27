import React, { useEffect, useState } from "react";
import api from "../api";

export default function Dashboard() {
  const [stats, setStats] = useState({
    users: 0,
    products: 0,
    orders: 0,
  });
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchStats = async () => {
      try {
        // You may need to create a backend endpoint like: GET /admin/dashboard
        // For now, we can call individual endpoints and count results
        const [usersRes, productsRes, ordersRes] = await Promise.all([
          api.get("/admin/users"),
          api.get("/products"),
          api.get("/admin/orders"),
        ]);

        // Filter to only include customers (exclude admins)
        const customers = usersRes.data.users.filter(
          (user) => user.role !== "admin"
        );

        setStats({
          users: customers.length,
          products: productsRes.data.products.length,
          orders: ordersRes.data.length,
        });

        // Store all customer data for the table
        setUsers(customers);
      } catch (err) {
        console.error(err);
        setError("Failed to load dashboard stats");
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) return <p className="p-4">Loading dashboard...</p>;
  if (error) return <p className="p-4 text-red-500">{error}</p>;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-6">Dashboard</h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-8">
        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-600">Total Users</h2>
          <p className="text-3xl font-bold text-blue-600 mt-2">{stats.users}</p>
        </div>

        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-600">Total Products</h2>
          <p className="text-3xl font-bold text-green-600 mt-2">
            {stats.products}
          </p>
        </div>

        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-600">Total Orders</h2>
          <p className="text-3xl font-bold text-purple-600 mt-2">
            {stats.orders}
          </p>
        </div>
      </div>

      {/* All Users List */}
      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-lg font-medium mb-4">All Users</h2>
        <table className="min-w-full border">
          <thead className="bg-gray-100">
            <tr>
              <th className="px-4 py-2 border">ID</th>
              <th className="px-4 py-2 border">Name</th>
              <th className="px-4 py-2 border">Email</th>
            </tr>
          </thead>
          <tbody>
            {users.length > 0 ? (
              users.map((user) => (
                <tr key={user.id}>
                  <td className="px-4 py-2 border">{user.id}</td>
                  <td className="px-4 py-2 border">{user.full_name}</td>
                  <td className="px-4 py-2 border">{user.email}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td
                  colSpan="3"
                  className="text-center py-4 text-gray-500 italic"
                >
                  No users found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
