import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie"; // импортируем js-cookie

import './styles/Auth.css';

const AuthPage = () => {
  const navigate = useNavigate();

  const [password, setPassword] = useState("");
  const [login, setLogin] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async () => {
    setError("");

    try {
      const response = await fetch("/api/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: login,
          password: password,
        }),
      });

      if (response.status === 200) {
        const data = await response.json();

        sessionStorage.setItem("access_token", data.access_token) // more security

        Cookies.set("refresh_token", data.refresh_token); // default 
        Cookies.set("token_type", data.token_type); // default

        // Перенаправляем на главную
        navigate("/");
      } else {
        const data = await response.json();
        setError(data.detail || "Login failed");
      }
    } catch (err) {
      setError("Something went wrong");
    }
  };

  const handleSignup = () => {
    navigate("/register");
  };

  return (
    <div className="auth-container">
      <div className="title">
        <h1>Welcome back to <br /> InnoQuiz</h1>
      </div>

      <div className="login-passwd-container">
        <div className="login-passwd-panel">
          <h1>Login</h1>
          <input
            type="text"
            placeholder="Email"
            value={login}
            onChange={(e) => setLogin(e.target.value)}
            className="login-passwd-input"
          />

          <h1>Password</h1>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="login-passwd-input"
          />

          {error && <div className="auth-error">{error}</div>}

          <button className="submit-button" onClick={handleSubmit}>Submit</button>

          <p className="signup-question">No account yet?</p>
          <button className="signup-button" onClick={handleSignup}>Sign Up</button>
        </div>
      </div>
    </div>
  );
};

export default AuthPage;
