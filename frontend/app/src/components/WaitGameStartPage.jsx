import React, {useEffect, useRef, useState} from "react";
import { useNavigate } from "react-router-dom";
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {

    const sessionServiceWsRef = useRef(null);
    const [players, setPlayers] = useState([]);

    const navigate = useNavigate();

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
                const incomingNames = JSON.parse(message.data); // пример: ["Alice"] или ["Alice", "Bob"]

                if (!Array.isArray(incomingNames)) return;

                setPlayers((prevPlayers) => {

                    // Фильтруем новых
                    const newPlayers = incomingNames
                        .map((name, index) => ({
                            id: prevPlayers.length + index + 1,
                            name: name
                        }));

                    return [...prevPlayers, ...newPlayers];
                });

                console.log("📨 Received JSON message:",incomingNames);
            } catch (e){
                console.error("⚠️ Failed to parse incoming WebSocket message:", message.data);
            }
        }
    };
    useEffect(() => {
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
