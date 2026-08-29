"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiRequest } from "../../../../services/api";

export default function LoginPage() {
  const router = useRouter();

  const [identify, setIdentify] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      await apiRequest("/api/login", "POST", { identify, password });
      router.push("/");

    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h2>Login Page</h2>

      {error && <p style={{ color: "red" }}>{error}</p>}

      <form onSubmit={handleSubmit}>
        <div>
          <label>Username or Email: </label>
          <input
            type="text"
            value={identify}
            onChange={(e) => setIdentify(e.target.value)}
            required
          />
        </div>

        <br />

        <div>
          <label>Password: </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <br />

        <button type="submit">Login</button>
      </form>
    </div>
  );
}