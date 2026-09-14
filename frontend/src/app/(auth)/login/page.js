"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { apiRequest } from "@/utils/api";
import { useAuth } from "@/app/context/AuthContext";

export default function LoginPage() {
  const router = useRouter();
  const { refreshAuth } = useAuth();

  const [formData, setFormData] = useState({ identify: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e) =>
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await apiRequest("/login", "POST", formData);
      await refreshAuth();
      router.push("/posts");
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  return (
    <div className="glass-page" style={{ flexDirection: 'column' }}>
      <div className="glass-container glass-container-login" style={{ minHeight: '500px', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
        <h1 style={{ fontSize: '32px', fontWeight: '800', marginBottom: '24px', marginTop: '0', color: '#fff', background: 'linear-gradient(to right, #fff, #a78bfa)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent', letterSpacing: '-0.02em', textAlign: 'center' }}>
          01-Social-Network
        </h1>
        <h2>Welcome Back</h2>
        <p>Sign in with your credentials</p>

        {error && <div className="glass-error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="glass-input-group">
            <label>Username / Email</label>
            <input
              type="text"
              className="glass-input"
              name="identify"
              required
              onChange={handleChange}
            />
          </div>

          <div className="glass-input-group">
            <label>Password</label>
            <input
              type="password"
              className="glass-input"
              name="password"
              required
              onChange={handleChange}
            />
          </div>

          <button type="submit" className="glass-button" disabled={loading}>
            {loading ? "Signing in..." : "Sign In"}
          </button>

          <div className="glass-footer">
            Don&apos;t have an account?
            <Link href="/register">Register</Link>
          </div>
        </form>
      </div>
    </div>
  );
}
