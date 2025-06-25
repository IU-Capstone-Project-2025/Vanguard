import React from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";

const PlayGamePage = () => {
    const [nickname,setNickname] = useState("")
    const navigate = useNavigate()


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

            const response = await fetch("/api/join-session", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ name: nickname, sessionCode }),
            });

            if (!response.ok) {
                throw new Error("Failed to join session.");
            }

            const data = await response.json();
            const { wsEndpoint, playerId } = data;

            sessionStorage.setItem("nickname", nickname);
            sessionStorage.setItem("WSEndpoint", wsEndpoint);
            sessionStorage.setItem("playerId", playerId); // если нужно

            navigate(`/wait/${sessionCode}`);
        } catch (error) {
            console.error("Join session error:", error);
            alert("Could not connect to the session. Please try again.");
        }
    };

    return (
        <div className="playgame-main-content">
            
                <div className="title">
                    <h1>
                        Now enter your nickname
                    </h1>
                    <input 
                        type="text" 
                        placeholder="enter the name here"
                        required
                        autoFocus
                        value={nickname}
                        onChange={(e)=> setNickname(e.target.value)}
                        className="code-input"
                    />
                    <div className="button-group">
                        <button id="play"
                                className="play-button"
                                onClick={
                                    (e) => {
                                        handlePlay();
                                        e.preventDefault();
                                    }
                                }
                            >
                            <span>Play</span>
                        </button>
                    </div>
                </div>
            </div>
    )
};

export default PlayGamePage;