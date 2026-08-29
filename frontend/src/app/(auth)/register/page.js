"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { apiRequest } from "../../../utils/api"; 

export default function RegisterPage() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    username: "",
    first_name: "",
    last_name: "",
    email: "",
    password: "",
    date_of_birth: "",
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
      await apiRequest("/api/register", "POST", formData);

      alert("Account created successfully!");
      router.push("/login");
      
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-96">
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-800">Create Account</h2>
        
        {error && <div className="bg-red-100 text-red-600 p-2 rounded mb-4 text-sm">{error}</div>}
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <input type="text" name="username" placeholder="Username" required onChange={handleChange}
            className="w-full p-2 border rounded text-black" />
            
          <div className="flex gap-2">
            <input type="text" name="first_name" placeholder="First Name" required onChange={handleChange}
              className="w-1/2 p-2 border rounded text-black" />
            <input type="text" name="last_name" placeholder="Last Name" required onChange={handleChange}
              className="w-1/2 p-2 border rounded text-black" />
          </div>

          <input type="email" name="email" placeholder="Email" required onChange={handleChange}
            className="w-full p-2 border rounded text-black" />
            
          <input type="password" name="password" placeholder="Password" required minLength="8" onChange={handleChange}
            className="w-full p-2 border rounded text-black" />
            
          <input type="date" name="date_of_birth" required onChange={handleChange}
            className="w-full p-2 border rounded text-gray-600" />

          <button type="submit" disabled={loading}
            className="w-full bg-blue-600 text-white p-2 rounded hover:bg-blue-700 transition disabled:bg-gray-400">
            {loading ? "Creating..." : "Register"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Already have an account? <Link href="/login" className="text-blue-600 hover:underline">Login here</Link>
        </p>
      </div>
    </div>
  );
}