import React, {useEffect, useRef} from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";

const JoinGamePage = () => {
    const [code,setCode] = useState("");
    const navigate = useNavigate()
    const wsRef = useRef(null);

    // 🎯 POST-запрос к /api/session/join
    const joinSession = async (quizId,userId) => {

        const response = await fetch("/api/session/join", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ quizId, userId}),
        });
        if (!response.ok) throw new Error("Failed to join session");

        const data = await response.json();
        return data; // возвращает объект вида: {"serverWsEndpoint": "string","jwt":"string", "sessionId":"string"}
    };

    useEffect(() => {
        console.log("code updated:", code);
    }, [code]);

    // 🌐 Устанавливаем WebSocket-соединение
    const connectToWebSocket = (wsEndpoint, token) => {
        wsRef.current = new WebSocket(`${wsEndpoint}?token=${token}`);
        wsRef.current.onopen = () => {
            console.log("✅ WebSocket connected");
        };

        wsRef.current.onerror = (err) => {
            console.error("❌ WebSocket error:", err);

        };
    };

    const handlePlay =  async () => {
        if (code) {
            const sessionData = await joinSession(code ,null)
            connectToWebSocket(sessionData.wsEndpoint,sessionData.jwt);
            sessionStorage.setItem('sessionCode', code); // Store the session code in session storage
            navigate(`/wait/${code}`); // Navigate to the waiting page with the session code
        }
    };

    return (
        <div className="joingame-main-content">
            <div className="left-side">
                <div className="title">
                    <h1>
                        Ask your quiz creator for a code
                    </h1>
                    <input 
                        type="text" 
                        placeholder="enter a code here"
                        value={`${code.toUpperCase()}`}
                        onChange={(e) => setCode(e.target.value)} // Remove '#' and convert to uppercase
                        required
                        autoFocus
                        pattern="^[A-Z0-9]+$" // Ensure only alphanumeric characters are
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
            <div className="right-side">
                
            </div>
        </div>
    )
};

export default JoinGamePage;