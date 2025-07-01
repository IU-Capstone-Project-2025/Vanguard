import React, {useEffect, useRef, useState} from "react";
import { useNavigate, useParams } from "react-router-dom";
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
    const navigate = useNavigate();
    const sessionCode = useParams();
    sessionStorage.setItem("sessionCode", sessionCode);
    
    
    const sessionServiceWsRef = useRef(null);
    const [players, setPlayers] = useState([
        { id: 1, name: "Alice" },
        { id: 2, name: "Bob" },
        { id: 3, name: "Charlie" },
        { id: 4, name: "Diana" },
        { id: 5, name: "Eva" },
        { id: 6, name: "Frank" },
        { id: 7, name: "Grace" },
        { id: 8, name: "Henry" },
        { id: 9, name: "Isabella" },
    ]);
    

    // 🌐 Устанавливаем WebSocket-соединение с Session Service
    const connectToWebSocket = (token) => {
        let serverWsEndpoint = "ws://localhost:8081/ws";
        sessionServiceWsRef.current = new WebSocket(`${serverWsEndpoint}?token=${token}`);
        sessionServiceWsRef.current.onopen = () => {
            console.log("✅ WebSocket connected to Session Service");
        }
        sessionServiceWsRef.current.onerror = (err) => {
            console.error("❌ WebSocket with Session Service error:", err);
        }
    }
    useEffect(() => {
        if (!sessionStorage.getItem("nickname")) {
            navigate("/enter-nickname");
            return;
        }
        connectToWebSocket(sessionStorage.getItem("jwt"));
    })
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
            {players.map((player) => (
                <div key={player.id} className="player-box">
                <span>{player.name}</span>
                </div>
            ))}
            </div>
        </div>
        </div>
    );
};

export default WaitGameStartPlayer;
