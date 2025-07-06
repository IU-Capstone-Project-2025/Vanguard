import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext"
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
  const navigate = useNavigate();
  const { sessionCode } = useParams();
  const { wsRefSession, connectSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime} = useRealtimeSocket();

  const [players, setPlayers] = useState([]);

  useEffect(() => {
    // Сохраняем sessionCode в sessionStorage (один раз)
    if (sessionCode) sessionStorage.setItem("sessionCode", sessionCode);

    const token = sessionStorage.getItem("jwt");
    const nickname = sessionStorage.getItem("nickname");

    if (!nickname) {
      navigate("/enter-nickname");
      return;
    }

    if (!wsRefSession.current || wsRefSession.current.readyState > 1) {
      connectSession(token, handleMessage);
    } else {
      wsRefSession.current.onmessage = handleMessage;
    }

    return () => {
      if (wsRefSession.current) {
        wsRefSession.current.onmessage = null;
      }
    };
  }, [connectSession, navigate, sessionCode, wsRefSession]);

  const handleMessage = (event) => {
    try {
      const incomingNames = JSON.parse(event.data); // ["Alice", "Bob", ...]

      if (!Array.isArray(incomingNames)) return;

      setPlayers((prevPlayers) => {
        const newNames = incomingNames.filter(
          (name) => !prevPlayers.some((p) => p.name === name)
        );

        const newPlayers = newNames.map((name, index) => ({
          id: prevPlayers.length + index + 1,
          name,
        }));

        return [...prevPlayers, ...newPlayers];
      });

      console.log("📨 Received player list:", incomingNames);
    } catch (err) {
      console.error("⚠️ Failed to parse WebSocket message:", event.data);
    }
  };

  const handleLeave = () => {
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
