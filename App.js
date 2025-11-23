import React, { useState } from "react";
import "./App.css";
import Login from "./Login";
import Items from "./Items";

export default function App() {
  const [token, setToken] = useState(null);

  return (
    <div className="app">
      <header className="app-header">
        <h1>
          <span style={{ fontWeight: 700, color: "#60a5fa" }}>RupeShop</span>{" "}
          <span style={{ fontWeight: 300, color: "#cbd5e1" }}>Store</span>
        </h1>

        {token && (
          <button className="btn btn-ghost" onClick={() => setToken(null)}>
            Logout
          </button>
        )}
      </header>

      <main className="app-main">
        {!token ? <Login setToken={setToken} /> : <Items token={token} />}
      </main>
    </div>
  );
}
