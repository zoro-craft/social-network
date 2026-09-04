"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
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
      setLoading(false);
      return;
    }
    if (avatar) {
      submitData.append("avatar", avatar);
    }

    try {
      await apiRequest("/register", "POST", submitData, true);
      alert("Account created successfully!");
      router.push("/login");
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  return (
    <div className="glass-page glass-page-register" style={{ flexDirection: 'column' }}>
      <div className="glass-container glass-container-register" style={{ display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
        <h1 style={{ fontSize: '32px', fontWeight: '800', marginBottom: '24px', marginTop: '0', color: '#fff', background: 'linear-gradient(to right, #fff, #a78bfa)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent', letterSpacing: '-0.02em', textAlign: 'center' }}>
          01-Social-Network
        </h1>
        <h2>Create Account</h2>
        <p>Join us today</p>

        {error && <div className="glass-error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="glass-row">
            <div className="glass-input-group">
              <label>First Name *</label>
              <input
                type="text"
                className="glass-input"
                name="first_name"
                placeholder="John"
                required
                maxLength={30}
                onChange={handleChange}
              />
            </div>
            <div className="glass-input-group">
              <label>Last Name *</label>
              <input
                type="text"
                className="glass-input"
                name="last_name"
                placeholder="Doe"
                required
                maxLength={30}
                onChange={handleChange}
              />
            </div>
          </div>

          <div className="glass-row">
            <div className="glass-input-group">
              <label>Email Address *</label>
              <input
                type="email"
                className="glass-input"
                name="email"
                placeholder="you@example.com"
                required
                maxLength={100}
                onChange={handleChange}
              />
            </div>
            <div className="glass-input-group">
              <label>Date of Birth *</label>
              <input
                type="date"
                className="glass-input"
                name="date_of_birth"
                required
                onChange={handleChange}
              />
            </div>
          </div>

          <div className="glass-row">
            <div className="glass-input-group">
              <label>Password *</label>
              <input
                type="password"
                className="glass-input"
                name="password"
                placeholder="••••••••"
                required
                minLength={8}
                maxLength={50}
                onChange={handleChange}
              />
            </div>
            <div className="glass-input-group">
              <label>Nickname (optional)</label>
              <input
                type="text"
                className="glass-input"
                name="nickname"
                placeholder="cooluser42"
                maxLength={30}
                onChange={handleChange}
              />
            </div>
          </div>

          <div className="glass-row glass-row-align-end">
            <div className="glass-input-group">
              <label>About Me (optional)</label>
              <textarea
                className="glass-textarea"
                name="about_me"
                rows="1"
                placeholder="Tell us about yourself..."
                maxLength={200}
                onChange={handleChange}
              />
            </div>
            <div className="glass-input-group">
              <label className="glass-file-label">
                Avatar / Image
                <input
                  className="glass-file-input"
                  type="file"
                  name="avatar"
                  accept="image/*"
                  onChange={handleFileChange}
                />
              </label>
            </div>
          </div>

          <button type="submit" className="glass-button" disabled={loading}>
            {loading ? "Creating Account..." : "Register"}
          </button>

          <div className="glass-footer">
            Already have an account?
            <Link href="/login">Log in here</Link>
          </div>
        </form>
      </div>
    </div>
  );
}