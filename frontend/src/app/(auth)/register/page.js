"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiRequest } from "../../../utils/api";

export default function RegisterPage() {
  const router = useRouter();

  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    email: "",
    password: "",
    date_of_birth: "",
    nickname: "",
    about_me: "",
  });
  
  const [avatar, setAvatar] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleFileChange = (e) => {
    setAvatar(e.target.files[0]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const submitData = new FormData();
    for (const key in formData) {
      submitData.append(key, formData[key]);
    }
    if (avatar && avatar.size > 10 * 1024 * 1024) {
      setError("Avatar must be less than 10MB");
      return;
    }
    if (avatar) {
      submitData.append("avatar", avatar);
    }

    try {
      await apiRequest("/api/register", "POST", submitData, true);
      alert("Account created successfully!");
      router.push("/login");
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-xl w-full max-w-4xl border border-gray-100">
        
        <div className="text-center mb-6">
          <h2 className="text-3xl font-bold text-gray-900">Create Account</h2>
          <p className="text-gray-500 mt-1 text-sm">Join us today!</p>
        </div>
        
        {error && <div className="bg-red-50 border-l-4 border-red-500 text-red-700 p-2 rounded mb-4 text-sm">{error}</div>}
        
        <form onSubmit={handleSubmit} className="space-y-6">
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            
            <div className="space-y-4">
              <h3 className="text-xs font-bold text-gray-400 uppercase tracking-widest border-b pb-2">Required Info</h3>
              
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">First Name *</label>
                  <input type="text" name="first_name" required onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-gray-700 mb-1">Last Name *</label>
                  <input type="text" name="last_name" required onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black" />
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">Email Address *</label>
                <input type="email" name="email" required onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black" />
              </div>
              
              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">Password *</label>
                <input type="password" name="password" required minLength="8" onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black" />
              </div>

              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">Date of Birth *</label>
                <input type="date" name="date_of_birth" required onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black" />
              </div>
            </div>

            <div className="space-y-4">
              <h3 className="text-xs font-bold text-gray-400 uppercase tracking-widest border-b pb-2">Profile (Optional)</h3>
              
              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">Nickname</label>
                <input type="text" name="nickname" onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black" />
              </div>
              
              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">About Me</label>
                <textarea name="about_me" rows="2" onChange={handleChange} className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-black resize-none"></textarea>
              </div>
              
              <div>
                <label className="block text-xs font-semibold text-gray-700 mb-1">Avatar / Image</label>
                <input type="file" name="avatar" accept="image/*" onChange={handleFileChange} className="w-full p-1.5 border border-gray-300 rounded-lg text-xs text-gray-600 file:mr-4 file:py-1.5 file:px-4 file:border-0 file:text-xs file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 cursor-pointer" />
              </div>
            </div>

          </div>

          <div className="pt-2">
            <button type="submit" disabled={loading} className="w-full bg-blue-600 text-white font-bold p-3 rounded-lg hover:bg-blue-700 transition duration-200 shadow disabled:bg-blue-400">
              {loading ? "Creating Account..." : "Register"}
            </button>
            <p className="mt-4 text-center text-sm text-gray-600">
              Already have an account? <a href="/login" className="text-blue-600 font-semibold hover:underline">Log in here</a>
            </p>
          </div>

        </form>
      </div>
    </div>
  );
}