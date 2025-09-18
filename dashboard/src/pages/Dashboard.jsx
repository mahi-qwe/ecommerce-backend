import React from "react";
import { useEffect, useState } from "react";
import api from "../api";

const Dashboard = () => {
  const [stats, setStats] = useState({
    users: 0,
    products: 0,
    orders: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [usersRes, productsRes, ordersRes] = await Promise.all([
          api.get("/users"),
          api.get("/products"),
          api.get("/orders"),
        ]);

        setStats({
          users: usersRes.data.length,
          products: productsRes.data.length,
          orders: ordersRes.data.length,
        });
      } catch (err) {
        setError("Failed to fetch dashboard data");
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) return <p className="p-4">Loading dashboard...</p>;
  if (error) return <p className="p-4 text-red-500">{error}</p>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-lg font-semibold">Total Users</h2>
          <p className="text-2xl mt-2">{stats.users}</p>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-lg font-semibold">Total Products</h2>
          <p className="text-2xl mt-2">{stats.products}</p>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-lg font-semibold">Total Orders</h2>
          <p className="text-2xl mt-2">{stats.orders}</p>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
