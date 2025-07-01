import React, {useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";

import './styles/Auth.css';

const AuthPage = () => {
    const navigate = useNavigate();

    const [password, setPassword] = useState("");
    const [login, setLogin] = useState("");


    const handleSubmit = () => {
        navigate("/")
    };
    const handleSignup = () => {
        navigate("/register")
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
                        placeholder="Value"
                        value={login}
                        onChange={(e) => setLogin(e.target.value)}
                        className="login-passwd-input"
                    />

                    <h1>Password</h1>
                    <input
                        type="password"
                        placeholder="Value"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="login-passwd-input"
                    />

                    <button className="submit-button" onClick={handleSubmit}>Submit</button>

                    <p className="signup-question">No have account yet?</p>
                    <button className="signup-button" onClick={handleSignup}>Sign Up</button>
                </div>
            </div>
        </div>
    );

};

export default AuthPage;
