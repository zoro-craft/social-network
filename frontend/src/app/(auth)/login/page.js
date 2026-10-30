"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { apiRequest } from "@/utils/api";
import { useAuth } from "@/app/context/AuthContext";

// Seeded demo accounts. Keep in sync with backend/internal/database/seed.go.
const DEMO_PASSWORD = "demo1234";
const DEMO_ACCOUNTS = [
  { nickname: "sara", name: "Sara Johnson" },
  { nickname: "ahmed", name: "Ahmed Ali" },
  { nickname: "nicole", name: "Nicole Martin" },
  { nickname: "alex", name: "Alex Brown" },
  { nickname: "lina", name: "Lina Haddad" },
  { nickname: "omar", name: "Omar Farouk" },
];

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

  const handleDemoLogin = async (nickname) => {
    setLoading(true);
    setError("");

    try {
      await apiRequest("/login", "POST", {
        identify: nickname,
        password: DEMO_PASSWORD,
      });
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
            Don't have an account?
            <Link href="/register">Register</Link>
          </div>
        </form>

        <div style={{ marginTop: 22, paddingTop: 18, borderTop: "1px solid var(--border)" }}>
          <p
            style={{
              margin: "0 0 10px",
              fontSize: 12,
              fontWeight: 600,
              letterSpacing: "0.04em",
              textTransform: "uppercase",
              color: "var(--text-muted)",
              textAlign: "center",
            }}
          >
            Or try a demo account
          </p>

          <div style={{ display: "grid", gridTemplateColumns: "repeat(2, minmax(0, 1fr))", gap: 8 }}>
            {DEMO_ACCOUNTS.map((account) => (
              <button
                key={account.nickname}
                type="button"
                className="glass-btn"
                disabled={loading}
                onClick={() => handleDemoLogin(account.nickname)}
                style={{ width: "100%", textAlign: "left", overflow: "hidden" }}
              >
                <strong
                  style={{
                    display: "block",
                    fontSize: 13,
                    fontWeight: 600,
                    color: "var(--text-primary)",
                    whiteSpace: "nowrap",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                  }}
                >
                  {account.name}
                </strong>
                <span style={{ fontSize: 11, color: "var(--text-muted)" }}>
                  @{account.nickname}
                </span>
              </button>
            ))}
          </div>

          <p style={{ margin: "10px 0 0", fontSize: 11, color: "var(--text-muted)", textAlign: "center" }}>
            Password for every demo account: {DEMO_PASSWORD}
          </p>
        </div>
      </div>
    </div>
  );
}
