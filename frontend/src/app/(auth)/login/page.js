"use client";

import { useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { apiRequest } from "@/utils/api";
import { useAuth } from "@/app/context/AuthContext";

// Demo accounts seeded by backend seeder, keep in sync
const DEMO_USERS = [
  {
    nickname: "yasmine",
    name: "Yasmine Haddad",
    avatarSrc: "/demo-avatars/yasmine.jpeg",
  },
  {
    nickname: "karim",
    name: "Karim Benali",
    avatarSrc: "/demo-avatars/karim.jpeg",
  },
  {
    nickname: "sofia",
    name: "Sofia Marchetti",
    avatarSrc: "/demo-avatars/sofia.jpeg",
  },
  {
    nickname: "lucas",
    name: "Lucas Moreau",
    avatarSrc: "/demo-avatars/lucas.jpeg",
  },
  {
    nickname: "ines",
    name: "Ines Laurent",
    avatarSrc: "/demo-avatars/ines.jpeg",
  },
  {
    nickname: "omar",
    name: "Omar Diallo",
    avatarSrc: "/demo-avatars/omar.jpeg",
  },
];

// keep in sync with backend seeddata.DemoPassword
const DEMO_PASSWORD = "Password123";

export default function LoginPage() {
  const router = useRouter();
  const { refreshAuth } = useAuth();

  const [formData, setFormData] = useState({ identify: "", password: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [demoLoading, setDemoLoading] = useState(null);

  const showDemo = process.env.NEXT_PUBLIC_SHOW_DEMO !== "false";
  const busy = loading || Boolean(demoLoading);

  const handleChange = (e) =>
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (busy) return;

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
    if (busy) return;

    setDemoLoading(nickname);
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
      setDemoLoading(null);
    }
  };

  return (
    <div className="glass-page" style={{ flexDirection: "column" }}>
      <div
        className="glass-container glass-container-login"
        style={{
          minHeight: "500px",
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
        }}
      >
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

          <button type="submit" className="glass-button" disabled={busy}>
            {loading ? "Signing in..." : "Sign In"}
          </button>

          <div className="glass-footer">
            Don&apos;t have an account?
            <Link href="/register">Register</Link>
          </div>
        </form>

        {showDemo && (
          <>
            <div className="demo-divider">
              <span>or try a demo account</span>
            </div>

            <div className="demo-users">
              {DEMO_USERS.map((user) => (
                <button
                  key={user.nickname}
                  type="button"
                  className="demo-user"
                  onClick={() => handleDemoLogin(user.nickname)}
                  disabled={busy}
                  aria-label={`Sign in as ${user.name}`}
                >
                  <Image
                    src={user.avatarSrc}
                    alt=""
                    width={34}
                    height={34}
                    className="demo-avatar"
                    unoptimized
                  />
                  <span className="demo-meta">
                    <strong>{user.nickname}</strong>
                    <small>
                      {demoLoading === user.nickname
                        ? "Signing in..."
                        : user.name}
                    </small>
                  </span>
                </button>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
