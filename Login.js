import React, { useState } from "react";

export default function Login({ setToken }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function login() {
    try {
      const res = await fetch("http://localhost:8080/users/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      if (!res.ok) {
        alert("Invalid username/password");
        return;
      }

      const data = await res.json();
      setToken(data.token);
    } catch (err) {
      alert("Login failed!");
    }
  }

  return (
    <div className="auth-card">
      <h2>Welcome back 👋</h2>
      <p className="auth-subtitle">Sign in to start shopping</p>

      <div className="form-group">
        <label>Username</label>
        <input
          className="input"
          placeholder="e.g. rupesh"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
      </div>

      <div className="form-group">
        <label>Password</label>
        <input
          className="input"
          type="password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
      </div>

      <button className="btn btn-primary" onClick={login} style={{ width: "100%", marginTop: 6 }}>
        Login
      </button>
    </div>
  );
}
