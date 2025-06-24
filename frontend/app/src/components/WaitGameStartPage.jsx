import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
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

    const navigate = useNavigate();

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
