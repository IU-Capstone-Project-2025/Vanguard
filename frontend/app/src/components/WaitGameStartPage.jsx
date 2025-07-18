import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext"
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
  const navigate = useNavigate();
  const { sessionCode } = useParams();
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();

  const [players, setPlayers] = useState(new Map());

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
      connectSession(token, handleMessageSession);
    } else {
      wsRefSession.current.onmessage = handleMessageSession;
    }

    if (!wsRefRealtime.current || wsRefRealtime.current.readyState > 1) {
      connectRealtime(token, handleMessageRealtime);
    } else {
      wsRefRealtime.current.onmessage = handleMessageRealtime;
    }

    wsRefRealtime.current.onclose = () => {
      endSession();
    }
    wsRefSession.current.onclose = () => {
      endSession();
    }

    return () => {
      if (wsRefSession.current) {
        wsRefSession.current.onmessage = null;
      }
      if (wsRefRealtime.current) {
        wsRefRealtime.current.onmessage = null;
      }
    };
  }, [connectSession, navigate, sessionCode, wsRefSession, wsRefRealtime, connectRealtime]);

  const endSession = () => {
    console.log(`Ending session... ${sessionCode}`);
    sessionStorage.removeItem('sessionCode');
    sessionStorage.removeItem('nickname')
    closeWsRefRealtime();
    closeWsRefSession();
    navigate('/');
  }


  const handleStartGame = async () => {
    if (!sessionCode) {
      console.error("Session code is not available");
      return;
    }

    navigate(`/game-process/${sessionCode}`);
  }

  const handleMessageRealtime = (event) => {
    try {
      const incomingData = JSON.parse(event.data);
      if (incomingData.type === "next_question") {
        console.log("📨 Received start game signal:", incomingData);
        handleStartGame();
      } else if (incomingData.type === "end_session") {
        console.log("📨 Received end session signal:", incomingData);
        endSession();
      } else {
        console.warn("⚠️ Unknown message type:", incomingData.type);
      }
    } catch (error) {
      console.error("Failed to parse WebSocket message:", event.data);
    }
  };

  const handleMessageSession = (event) => {
    try {
      const data = JSON.parse(event.data); // ["Alice", "Bob", ...]

      setPlayers(() => {
        const newPlayers = new Map();
        for (const [userId,name] of Object.entries(data)) {
          if (!newPlayers.has(userId)) {
            newPlayers.set(userId, name);
          }
        }
        return newPlayers;
      });

      console.log("📨 Received player list:", data);
    } catch (err) {
      console.error("⚠️ Failed to parse WebSocket message:", event.data);
    }
  };

  const handleLeave = () => {
    endSession();
  };

  return (
    <div className="wait-player-container">
      <div className="wait-player-panel">
        <h1>Now let's wait for your friends</h1>
        <div className="button-group">
          <button onClick={handleLeave}>🔙 Leave</button>
        </div>
        <div className="players-grid">
          {Array.from(players.entries()).map(([id,name]) => (
            name !== "Admin" && 
            <div key={id} className="player-box">
              <span style={{ '--name-length': +name.length }}>{name}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WaitGameStartPlayer;
