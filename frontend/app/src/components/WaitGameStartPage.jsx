import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
  const [players, setPlayers] = useState([]);
  const ws = useRef(null);
  const navigate = useNavigate();

  useEffect(() => {
    const wsEndpoint = sessionStorage.getItem("WSEndpoint");
    if (!wsEndpoint) {
      console.error("WebSocket endpoint not found in sessionStorage.");
      return;
    }

    ws.current = new WebSocket(wsEndpoint);

    ws.current.onopen = () => {
      console.log("Player WebSocket connected.");
      ws.current.send(JSON.stringify({ type: "playerReady" }));
    };

    ws.current.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        switch (message.type) {
          case "newPlayerJoined":
            setPlayers((prev) => {
              const alreadyExists = prev.some((p) => p.id === message.player.id);
              return alreadyExists ? prev : [...prev, message.player];
            });
            break;

          case "currentPlayersList":
            setPlayers(message.players);
            break;

          case "playerKicked":
            setPlayers((prev) => prev.filter((p) => p.id !== message.playerId));
            break;

          case "startGame":
            const sessionCode = sessionStorage.getItem("sessionCode");
            navigate(`/game-player/${sessionCode}`);
            break;

          case "sessionTerminated":
            alert("Сессия была завершена администратором.");
            navigate("/");
            break;

          default:
            console.warn("Unknown message type:", message.type);
        }
      } catch (e) {
        console.error("WebSocket message parse error:", e);
      }
    };

    ws.current.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    ws.current.onclose = () => {
      console.log("WebSocket connection closed.");
    };

    return () => {
      ws.current?.close();
    };
  }, [navigate]);

  const handleLeave = () => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: "playerLeave" }));
    }
    navigate("/");
  };

  return (
    <div className="wait-player-container">
      <div className="wait-player-panel">
        <h1>Now let's wait for your friends</h1>
        <div className="button-group">
          <button onClick={handleLeave}>🔙 Leave</button>
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
