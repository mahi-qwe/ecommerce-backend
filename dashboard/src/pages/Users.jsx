import React, { useEffect, useState } from "react";
import api from "../api";
import { MdDelete, MdModeEdit, MdBlock } from "react-icons/md";
import { CgUnblock } from "react-icons/cg";

export default function Users() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [isEditing, setIsEditing] = useState(false);
  const [currentId, setCurrentId] = useState(null);
  const [showForm, setShowForm] = useState(false);

  const [form, setForm] = useState({
    full_name: "",
    role: "",
    address: "",
    avatar_url: "",
  });

  // Fetch users
  const fetchUsers = async () => {
    try {
      const res = await api.get("/admin/users");
      setUsers(res.data.users);
    } catch (err) {
      console.error(err);
      setError("Failed to load users");
    } finally {
      setLoading(false);
    }
  };

  // Reset form
  const resetForm = () => {
    setForm({
      full_name: "",
      role: "",
      address: "",
      avatar_url: "",
    });
    setIsEditing(false);
    setCurrentId(null);
    setShowForm(false);
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  // Handle input change
  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm({ ...form, [name]: value });
  };

  // Handle update user
  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const payload = {
        full_name: form.full_name,
        role: form.role,
        address: form.address,
        avatar_url: form.avatar_url,
      };

      await api.put(`/admin/users/${currentId}`, payload);
      alert("✅ User updated successfully");
      resetForm();
      fetchUsers();
    } catch (err) {
      console.error("Update user error:", err);
      alert("❌ Failed to update user");
    }
  };

  // Handle edit
  const handleEdit = (user) => {
    setForm({
      full_name: user.full_name,
      role: user.role,
      address: user.address || "",
      avatar_url: user.avatar_url || "",
    });
    setIsEditing(true);
    setCurrentId(user.id);
    setShowForm(true);
  };

  // Handle delete
  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this user?")) return;
    try {
      await api.delete(`/admin/users/${id}`);
      setUsers(users.filter((u) => u.id !== id));
    } catch (err) {
      console.error(err);
      alert("Failed to delete user");
    }
  };

  // Handle block
  const handleBlock = async (id) => {
    try {
      await api.post(`/admin/users/${id}/block`);
      setUsers(
        users.map((u) => (u.id === id ? { ...u, is_blocked: true } : u))
      );
    } catch (err) {
      console.error(err);
      alert("Failed to block user");
    }
  };

  // Handle unblock
  const handleUnblock = async (id) => {
    try {
      await api.post(`/admin/users/${id}/unblock`);
      setUsers(
        users.map((u) => (u.id === id ? { ...u, is_blocked: false } : u))
      );
    } catch (err) {
      console.error(err);
      alert("Failed to unblock user");
    }
  };

  if (loading) return <p className="p-4">Loading...</p>;
  if (error) return <p className="p-4 text-red-500">{error}</p>;

  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold">Users</h1>
      </div>

      {/* Edit Form - only shown when showForm is true */}
      {showForm && (
        <form
          onSubmit={handleSubmit}
          className="mb-6 space-y-4 bg-white p-4 shadow rounded"
        >
          <input
            type="text"
            name="full_name"
            value={form.full_name}
            onChange={handleChange}
            placeholder="Full Name"
            className="w-full p-2 border rounded"
            required
          />
          <select
            name="role"
            value={form.role}
            onChange={handleChange}
            className="w-full p-2 border rounded"
            required
          >
            <option value="">Select Role</option>
            <option value="user">User</option>
            <option value="admin">Admin</option>
          </select>
          <textarea
            name="address"
            value={form.address}
            onChange={handleChange}
            placeholder="Address"
            className="w-full p-2 border rounded"
          />
          <input
            type="text"
            name="avatar_url"
            value={form.avatar_url}
            onChange={handleChange}
            placeholder="Avatar URL"
            className="w-full p-2 border rounded"
          />
          {/* Preview */}
          {form.avatar_url && (
            <img
              src={form.avatar_url}
              alt="Preview"
              className="w-32 h-32 object-cover rounded border"
            />
          )}
          <div className="space-x-2">
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Update User
            </button>
            <button
              type="button"
              onClick={resetForm}
              className="px-4 py-2 bg-gray-400 text-white rounded hover:bg-gray-500"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {/* Users Table */}
      <table className="min-w-full border border-gray-300 rounded-lg overflow-hidden bg-white shadow">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-4 py-2 border">ID</th>
            <th className="px-4 py-2 border">Avatar</th>
            <th className="px-4 py-2 border">Name</th>
            <th className="px-4 py-2 border">Email</th>
            <th className="px-4 py-2 border">Role</th>
            <th className="px-4 py-2 border">Address</th>
            <th className="px-4 py-2 border">Status</th>
            <th className="px-4 py-2 border">Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.filter((u) => u.role === "user").length > 0 ? (
            users
              .filter((u) => u.role === "user")
              .map((u) => (
                <tr key={u.id} className="text-center">
                  <td className="px-4 py-2 border">{u.id}</td>
                  <td className="px-4 py-2 border">
                    {u.avatar_url ? (
                      <img
                        src={u.avatar_url}
                        alt={u.full_name}
                        className="w-16 h-16 object-cover rounded mx-auto"
                      />
                    ) : (
                      <span className="text-gray-400 italic">No avatar</span>
                    )}
                  </td>
                  <td className="px-4 py-2 border">{u.full_name}</td>
                  <td className="px-4 py-2 border">{u.email}</td>
                  <td className="px-4 py-2 border">{u.role}</td>
                  <td className="px-4 py-2 border">{u.address || "N/A"}</td>
                  <td className="px-4 py-2 border">
                    {u.is_blocked ? (
                      <span className="text-red-500 font-semibold">
                        Blocked
                      </span>
                    ) : (
                      <span className="text-green-600 font-semibold">
                        Active
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-2 border space-x-2">
                    <button
                      onClick={() => handleEdit(u)}
                      className="px-1 py-1 bg-gray-300 rounded-full"
                    >
                      <MdModeEdit />
                    </button>
                    {u.is_blocked ? (
                      <button
                        onClick={() => handleUnblock(u.id)}
                        className="px-1 py-1 bg-gray-300 rounded-full"
                      >
                        <CgUnblock />
                      </button>
                    ) : (
                      <button
                        onClick={() => handleBlock(u.id)}
                        className="px-1 py-1 bg-gray-300 rounded-full"
                      >
                        <MdBlock />
                      </button>
                    )}
                    <button
                      onClick={() => handleDelete(u.id)}
                      className="px-1 py-1 bg-gray-300 rounded-full"
                    >
                      <MdDelete />
                    </button>
                  </td>
                </tr>
              ))
          ) : (
            <tr>
              <td colSpan="8" className="text-center py-4 text-gray-500 italic">
                No customers found
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
