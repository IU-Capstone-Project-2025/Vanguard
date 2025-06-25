import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { v4 as uuidv4 } from "uuid"; // ✅ импорт
import './styles/styles.css';

const PlayGamePage = () => {
    const [nickname, setNickname] = useState("");
    const navigate = useNavigate();

    const handlePlay = async () => {
        if (!nickname.trim()) {
            alert("Please enter a nickname to continue.");
            return;
        }

        try {
            const sessionCode = sessionStorage.getItem("sessionCode");
            if (!sessionCode) {
                alert("Session code is missing. Please start from the beginning.");
                return;
            }

            // ✅ генерируем уникальный userId
            const userId = `${nickname}_${uuidv4()}`;

            const response = await fetch("/join", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    userId,      // ✅ добавляем userId в запрос
                    code: sessionCode,
                }),
            });

            if (!response.ok) {
                throw new Error("Failed to join session.");
            }

            const data = await response.json();
            const { jwt, serverWsEndpoint, sessionId } = data;

            // ✅ сохраняем userId и nickname
            sessionStorage.setItem("nickname", nickname);
            sessionStorage.setItem("userId", userId);
            sessionStorage.setItem("WSEndpoint", serverWsEndpoint);
            sessionStorage.setItem("JWT", jwt);

            navigate(`/wait/${sessionCode}`);
        } catch (error) {
            console.error("Join session error:", error);
            alert("Could not connect to the session. Please try again.");
        }
    };

    return (
        <div className="playgame-main-content">
            <div className="title">
                <h1>Now enter your nickname</h1>
                <input 
                    type="text" 
                    placeholder="enter the name here"
                    required
                    autoFocus
                    value={nickname}
                    onChange={(e) => setNickname(e.target.value)}
                    className="code-input"
                />
                <div className="button-group">
                    <button
                        id="play"
                        className="play-button"
                        onClick={(e) => {
                            e.preventDefault();
                            handlePlay();
                        }}
                    >
                        <span>Play</span>
                    </button>
                </div>
            </div>
        </div>
    );
};

export default PlayGamePage;
