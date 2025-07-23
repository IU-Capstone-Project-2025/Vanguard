import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { API_ENDPOINTS } from '../constants/api';

import styles from './styles/Auth.module.css';

const AuthPage = () => {
  const navigate = useNavigate();

  const [password, setPassword] = useState("");
  const [login, setLogin] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async () => {
    setError("");

    try {
      const response = await fetch(`${API_ENDPOINTS.AUTH}/login`, {
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

        sessionStorage.setItem("access_token", data.access_token);

        Cookies.set("refresh_token", data.refresh_token);
        Cookies.set("token_type", data.token_type);

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
    <div className={styles["auth-container"]}>
      <div className={styles["login-passwd-container"]}>
        <div className={styles["login-passwd-panel"]}>
          <h1>Login</h1>
          <input
            type="text"
            placeholder="Email"
            value={login}
            onChange={(e) => setLogin(e.target.value)}
            className={styles["login-passwd-input"]}
          />

          <h1>Password</h1>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={styles["login-passwd-input"]}
          />

          {error && <div className={styles["auth-error"]}>{error}</div>}

          <button className={styles["submit-button"]} onClick={handleSubmit}>Submit</button>

          <p className={styles["signup-question"]}>No account yet?</p>
          <button className={styles["signup-button"]} onClick={handleSignup}>Sign Up</button>
        </div>
      </div>
    </div>
  );
}

export default AuthPage