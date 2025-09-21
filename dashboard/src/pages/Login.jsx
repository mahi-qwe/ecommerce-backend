import React from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api";

// Helper function to decode JWT and extract payload
const decodeJWT = (token) => {
  try {
    const payload = token.split(".")[1];
    const decoded = JSON.parse(atob(payload));
    return decoded;
  } catch (error) {
    console.error("Error decoding JWT:", error);
    return null;
  }
};

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await api.post("/auth/login", { email, password });
      const token = res.data.access_token;

      // Decode JWT to check user role
      const decoded = decodeJWT(token);

      if (!decoded) {
        setError("Invalid token received");
        setLoading(false);
        return;
      }

      // Check if user is admin
      if (decoded.role !== "admin") {
        setError("Access denied. Admin privileges required.");
        setLoading(false);
        return;
      }

      // Save JWT token in localStorage only if user is admin
      localStorage.setItem("token", token);
      localStorage.setItem("userRole", decoded.role);
      localStorage.setItem("userId", decoded.user_id);

      // Redirect to dashboard
      navigate("/dashboard");
    } catch (err) {
      setError(err.response?.data?.error || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md bg-white p-8 rounded-2xl shadow-lg">
        <h2 className="text-2xl font-bold mb-6 text-center">Admin Login</h2>

        {error && (
          <p className="text-red-500 text-sm mb-4 text-center">{error}</p>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-600">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={loading}
              className="mt-1 w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring focus:ring-blue-300 disabled:opacity-50"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-600">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={loading}
              className="mt-1 w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring focus:ring-blue-300 disabled:opacity-50"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <p className="text-xs text-gray-500 text-center mt-4">
          This portal is restricted to administrators only
        </p>
      </div>
    </div>
  );
}
