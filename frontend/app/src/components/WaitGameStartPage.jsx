import React, {useEffect, useRef, useState} from "react";
import { useNavigate, useParams } from "react-router-dom";
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
    const navigate = useNavigate();
    const sessionCode = useParams();
    sessionStorage.setItem("sessionCode", sessionCode);
    
    
    const sessionServiceWsRef = useRef(null);

    const [players, setPlayers] = useState(new Map());

    // 🌐 Устанавливаем WebSocket-соединение с Session Service
    const connectToWebSocket = (token) => {
        let serverWsEndpoint = "ws://localhost:8081/ws";
        sessionServiceWsRef.current = new WebSocket(`${serverWsEndpoint}?token=${token}`);
        sessionServiceWsRef.current.onopen = () => {
            console.log("✅ WebSocket connected with Session Service");
        };

        sessionServiceWsRef.current.onerror = (err) => {
            console.error("❌ WebSocket with Session Service error:", err);

        };

        // получение сообщения от session service

        sessionServiceWsRef.current.onmessage = (message) => {
            try {
                const incoming = JSON.parse(message.data); // { user123: "Alice", user456: "Bob" }

                setPlayers((prevMap) => {
                    const updatedMap = new Map(prevMap); // Копируем старую Map

                    for (const [userId, name] of Object.entries(incoming)) {
                        if (!updatedMap.has(userId)) {
                            updatedMap.set(userId, name); // Добавляем только новых
                        }
                    }

                    return updatedMap;
                });

                console.log("📨 Received JSON message:",incoming);
            } catch (e){
                console.error("⚠️ Failed to parse incoming WebSocket message:", message.data);
            }
        }
    };
    useEffect(() => {
        if (!sessionStorage.getItem("nickname")) {
            navigate("/enter-nickname");
            return;
        }
        connectToWebSocket(sessionStorage.getItem("jwt"));
    },[])



    const handleLeave = () => {
        navigate("/");
    };

    return (
        <div className="wait-player-container">
        <div className="wait-player-panel">
            <h1>Now let's wait for your friends</h1>
            <div className="button-group">
            <button onClick={handleLeave}>🔙Leave</button>
            </div>
            <div className="players-grid">
                {Array.from(players.entries()).map(([id, name]) => (
                    <div key={id} className="player-box">
                        <span>{name}</span>
                </div>
            ))}
            </div>
        </div>
        </div>
    );
};

export default WaitGameStartPlayer;
