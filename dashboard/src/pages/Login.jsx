import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api";
import { Lock, Shield } from "lucide-react";

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

      localStorage.setItem("token", token);
      localStorage.setItem("userRole", decoded.role);
      localStorage.setItem("userId", decoded.user_id);
      navigate("/dashboard");
    } catch (err) {
      setError(err.response?.data?.error || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-slate-800 via-slate-900 to-slate-950">
      <div className="w-full max-w-md p-8 rounded-2xl bg-slate-900 shadow-2xl border border-blue-600">
        {/* Heading */}
        <div className="flex justify-center items-center mb-6 gap-3">
          <div className="w-12 h-12 bg-blue-500 rounded-lg flex items-center justify-center shadow-lg">
            <Shield className="w-6 h-6 text-white" />
          </div>
          <h2 className="text-2xl font-extrabold text-white tracking-widest">
            Admin Login
          </h2>
        </div>

        {error && (
          <p className="text-red-500 text-xs font-semibold mb-4 text-center bg-slate-800 rounded-lg py-2 px-4 shadow">
            {error}
          </p>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-sm font-medium text-blue-200 mb-1">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={loading}
              className="mt-1 w-full px-4 py-2 bg-slate-800 border border-blue-500 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-400 disabled:opacity-50 placeholder-gray-400"
              placeholder="admin@email.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-blue-200 mb-1">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={loading}
              className="mt-1 w-full px-4 py-2 bg-slate-800 border border-blue-500 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-400 disabled:opacity-50 placeholder-gray-400"
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 px-5 py-2 rounded-lg font-bold text-white text-lg transition-all duration-200 shadow-lg transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Lock className="w-5 h-5 text-white" />
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <p className="text-xs text-blue-300 text-center mt-6">
          This portal is{" "}
          <span className="font-bold text-blue-500">restricted</span> to
          administrators only.
        </p>
      </div>
    </div>
  );
}
