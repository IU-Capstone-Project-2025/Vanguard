import React, { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import "./styles/WaitGameStartAdmin.css";

const WaitGameStartAdmin = () => {
  const navigate = useNavigate();
  const [players, setPlayers] = useState([]);
  const ws = useRef(null);

  // Подключаемся к WebSocket при монтировании компонента
  useEffect(() => {
    const wsEndpoint = sessionStorage.getItem("WSEndpoint");
    if (!wsEndpoint) {
      console.error("WebSocket endpoint not found in sessionStorage.");
      return;
    }

    ws.current = new WebSocket(wsEndpoint);

    ws.current.onopen = () => {
      console.log("WebSocket connection established.");
      // Можно послать init-сообщение или запросить уже подключившихся игроков
      ws.current.send(JSON.stringify({ type: "adminConnected" }));
    };

    ws.current.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        if (message.type === "newPlayerJoined") {
          const { id, name } = message.player;
          setPlayers((prev) => {
            const alreadyExists = prev.some((p) => p.id === id);
            return alreadyExists ? prev : [...prev, { id, name }];
          });
        }
        if (message.type === "currentPlayersList") {
          // В случае, если сервер отправит полный список игроков
          setPlayers(message.players);
        }
        if (message.type === "playerKicked") {
          setPlayers((prev) => prev.filter((p) => p.id !== message.playerId));
        }
      } catch (e) {
        console.error("Error parsing WebSocket message:", e);
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
  }, []);

  const handleKick = (idToRemove) => {
    setPlayers((prev) => prev.filter((player) => player.id !== idToRemove));
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: "kickPlayer", playerId: idToRemove }));
    }
  };

  const handleStart = () => {
    const sessionCode = sessionStorage.getItem("sessionCode");
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: "startGame" }));
    }
    navigate(`/game-controller/${sessionCode}`);
  };

  const handleTerminate = () => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ type: "terminateSession" }));
    }
    navigate("/");
  };

  return (
    <div className="wait-admin-container">
      <div className="wait-admin-panel">
        <h1>Now let's wait for your friends</h1>
        <div className="admin-button-group">
          <button onClick={handleStart}>▶ Start</button>
          <button onClick={handleTerminate}>▶ Terminate</button>
        </div>
        <div className="players-grid">
          {players.map((player) => (
            <div key={player.id} className="player-box">
              <span>{player.name}</span>
              <button
                className="kick-button"
                onClick={() => handleKick(player.id)}
                title={`Kick ${player.name}`}
              >
                ❌
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WaitGameStartAdmin;
