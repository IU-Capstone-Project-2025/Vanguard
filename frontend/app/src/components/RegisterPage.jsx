import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import styles from './styles/RegisterPage.module.css';
import { API_ENDPOINTS } from '../constants/api';

const RegisterPage = () => {
  const [formData, setFormData] = useState({
    email: "",
    nickname: "",
    password: "",
    repeatPassword: ""
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (formData.password !== formData.repeatPassword) {
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
          email: formData.email,
          username: formData.nickname,
          password: formData.password,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        Cookies.set("user_id", data.id);
        Cookies.set("user_email", data.email);
        Cookies.set("user_nickname", data.username);
        setSuccess("Registration successful! Redirecting...");
        setTimeout(() => navigate("/login"), 1500);
      } else {
        const errorData = await response.json();
        setError(errorData.message || "Registration failed");
      }
    } catch (err) {
      setError("Network error. Please try again.");
    }
  };

  return (
    <div className={styles['register-container']}>
      <div className={styles['register-box']}>
        <h2 className={styles['register-title']}>Sign up to InnoQuiz</h2>

        <form onSubmit={handleSubmit} className={styles['register-form']}>
          <input
            type="email"
            name="email"
            placeholder="Email"
            value={formData.email}
            required
            onChange={handleChange}
            className={styles['register-input']}
          />
          <input
            type="text"
            name="nickname"
            placeholder="Nickname"
            value={formData.nickname}
            required
            onChange={handleChange}
            className={styles['register-input']}
          />
          <input
            type="password"
            name="password"
            placeholder="Password"
            value={formData.password}
            required
            onChange={handleChange}
            className={styles['register-input']}
          />
          <input
            type="password"
            name="repeatPassword"
            placeholder="Repeat Password"
            value={formData.repeatPassword}
            required
            onChange={handleChange}
            className={styles['register-input']}
          />

          {error && <div className={styles['register-error']}>{error}</div>}
          {success && <div className={styles['register-success']}>{success}</div>}

          <button type="submit" className={styles['register-button']}>
            Register
          </button>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;