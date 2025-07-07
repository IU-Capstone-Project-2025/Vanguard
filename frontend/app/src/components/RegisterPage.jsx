import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie"; // добавляем библиотеку для работы с cookies
import "./styles/RegisterPage.css";
import { API_ENDPOINTS } from '../constants/api';

const RegisterPage = () => {
  const [email, setEmail] = useState("");
  const [nickname, setNickname] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const navigate = useNavigate(); // используем хук useNavigate для перенаправления

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (password !== repeatPassword) {
      setError("Passwords do not match");
      return;
    }

    try {
      const response = await fetch(`${API_ENDPOINTS.AUTH}/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          username: nickname,
          password: password,
        }),
      });

      if (response.status === 201) {
        const data = await response.json();

        // 💾 сохраняем данные в cookies
        Cookies.set("user_id", data.id);
        Cookies.set("user_email", data.email);
        Cookies.set("user_nickname", data.username);

        setSuccess("Registration successful! You can now log in.");
        navigate("/login"); // перенаправляем на страницу логина
      } else {
        const data = await response.json();
        setError(data.detail?.[0]?.msg || "Registration failed");
      }
    } catch (err) {
      setError("Something went wrong");
    }
  };

  return (
    <div className="register-container">
      <div className="register-box">
        <h2 className="register-title">Sign up to InnoQuiz</h2>

        <form onSubmit={handleSubmit} className="register-form">
          <input
            type="email"
            placeholder="Email"
            value={email}
            required
            onChange={(e) => setEmail(e.target.value)}
            className="register-input"
          />
          <input
            type="text"
            placeholder="Nickname"
            value={nickname}
            required
            onChange={(e) => setNickname(e.target.value)}
            className="register-input"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            required
            onChange={(e) => setPassword(e.target.value)}
            className="register-input"
          />
          <input
            type="password"
            placeholder="Repeat Password"
            value={repeatPassword}
            required
            onChange={(e) => setRepeatPassword(e.target.value)}
            className="register-input"
          />

          {error && <div className="register-error">{error}</div>}
          {success && <div className="register-success">{success}</div>}

          <button type="submit" className="register-button">Register</button>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;
