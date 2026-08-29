"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { apiRequest } from "../../../utils/api";

export default function LoginPage() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    identify: "", 
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true); 
    setError("");

    try {
      await apiRequest("/api/login", "POST", formData);

      alert("Logged in successfully!");
      router.push("/");
      
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-96">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">Login</h2>
        
        {error && <div className="bg-red-100 text-red-600 p-2 rounded mb-4 text-sm">{error}</div>}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <input 
            type="text" 
            name="identify" 
            placeholder="Username or Email" 
            required 
            onChange={handleChange}
            className="w-full p-2 border rounded text-black" 
          />
            
          <input 
            type="password" 
            name="password" 
            placeholder="Password" 
            required 
            onChange={handleChange}
            className="w-full p-2 border rounded text-black" 
          />

          <button 
            type="submit" 
            disabled={loading}
            className="w-full bg-green-600 text-white p-2 rounded hover:bg-green-700 transition disabled:bg-gray-400">
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Don't have an account? <Link href="/register" className="text-blue-600 hover:underline">Register here</Link>
        </p>
      </div>
    </div>
  );
}
