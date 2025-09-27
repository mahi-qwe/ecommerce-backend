import React, { useEffect, useState } from "react";
import api from "../api";

export default function Orders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filterStatus, setFilterStatus] = useState("all");

  const statusOptions = ["pending", "processing", "shipped", "delivered"];

  // Fetch all orders
  const fetchOrders = async () => {
    try {
      const res = await api.get("/admin/orders");
      setOrders(res.data);
    } catch (err) {
      console.error(err);
      setError("Failed to load orders");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders();
  }, []);

  // Handle status update
  const handleStatusChange = async (id, newStatus) => {
    if (!window.confirm(`Update order #${id} to ${newStatus}?`)) return;

    try {
      const res = await api.put(`/admin/orders/${id}`, { status: newStatus });

      // Update local state
      setOrders(
        orders.map((o) => (o.id === id ? { ...o, status: res.data.status } : o))
      );

      alert("✅ Order updated successfully");
    } catch (err) {
      console.error("Update status error:", err);
      alert("❌ Failed to update order");
    }
  };

  // Apply filtering
  const filteredOrders =
    filterStatus === "all"
      ? orders
      : orders.filter((o) => o.status === filterStatus);

  if (loading) return <p className="p-4">Loading...</p>;
  if (error) return <p className="p-4 text-red-500">{error}</p>;

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold">Orders</h1>

        {/* Filter Dropdown */}
        <select
          value={filterStatus}
          onChange={(e) => setFilterStatus(e.target.value)}
          className="p-2 border rounded"
        >
          <option value="all">All</option>
          {statusOptions.map((status) => (
            <option key={status} value={status}>
              {status.charAt(0).toUpperCase() + status.slice(1)}
            </option>
          ))}
        </select>
      </div>

      {/* Orders Table */}
      <table className="min-w-full border border-gray-300 rounded-lg overflow-hidden bg-white shadow">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-4 py-2 border">ID</th>
            <th className="px-4 py-2 border">Customer</th>
            <th className="px-4 py-2 border">Address</th>
            <th className="px-4 py-2 border">Total</th>
            <th className="px-4 py-2 border">Status</th>
            <th className="px-4 py-2 border">Items</th>
            <th className="px-4 py-2 border">Created</th>
            <th className="px-4 py-2 border">Actions</th>
          </tr>
        </thead>
        <tbody>
          {filteredOrders.length > 0 ? (
            filteredOrders.map((o) => (
              <tr key={o.id} className="text-center">
                <td className="px-4 py-2 border">{o.id}</td>
                <td className="px-4 py-2 border">{o.user_name}</td>
                <td className="px-4 py-2 border">{o.address}</td>
                <td className="px-4 py-2 border">₹{o.total_amount}</td>
                <td className="px-4 py-2 border">
                  <span
                    className={`px-2 py-1 rounded text-sm font-semibold ${
                      o.status === "delivered"
                        ? "bg-green-200 text-green-700"
                        : o.status === "shipped"
                        ? "bg-blue-200 text-blue-700"
                        : o.status === "processing"
                        ? "bg-yellow-200 text-yellow-700"
                        : "bg-gray-200 text-gray-700"
                    }`}
                  >
                    {o.status}
                  </span>
                </td>
                <td className="px-4 py-2 border text-left">
                  <ul className="list-disc list-inside text-sm">
                    {o.items.map((item, idx) => (
                      <li key={idx}>
                        {item.name} × {item.quantity} (₹{item.price})
                      </li>
                    ))}
                  </ul>
                </td>
                <td className="px-4 py-2 border">
                  {new Date(o.created_at).toLocaleString()}
                </td>
                <td className="px-4 py-2 border">
                  <select
                    value={o.status}
                    onChange={(e) => handleStatusChange(o.id, e.target.value)}
                    className="p-1 border rounded"
                  >
                    {statusOptions.map((status) => (
                      <option key={status} value={status}>
                        {status}
                      </option>
                    ))}
                  </select>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="8" className="text-center py-4 text-gray-500 italic">
                No orders found
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
