import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext";
import "./styles/WaitGameStartPlayer.css";

const WaitGameStartPlayer = () => {
  const navigate = useNavigate();
  const { sessionCode } = useParams();
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();

  const [players, setPlayers] = useState(new Map());

  useEffect(() => {
    const token = sessionStorage.getItem("jwt");
    const nickname = sessionStorage.getItem("nickname");
    if (!nickname) {
      navigate("/enter-nickname");
      return;
    }

    if (sessionCode) sessionStorage.setItem("sessionCode", sessionCode);

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

    wsRefRealtime.current.onclose = endSession;
    wsRefSession.current.onclose = endSession;

    return () => {
      if (wsRefSession.current) wsRefSession.current.onmessage = null;
      if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
    };
  }, [connectSession, connectRealtime, navigate, sessionCode, wsRefSession, wsRefRealtime]);

  const endSession = () => {
    sessionStorage.removeItem("sessionCode");
    sessionStorage.removeItem("nickname");
    closeWsRefRealtime();
    closeWsRefSession();
    navigate("/");
  };

  const handleStartGame = () => {
    if (sessionCode) navigate(`/game-process/${sessionCode}`);
  };

  const handleMessageRealtime = (event) => {
    try {
      const incomingData = JSON.parse(event.data);
      if (incomingData.type === "next_question") {
        handleStartGame();
      } else if (incomingData.type === "end_session") {
        endSession();
      }
    } catch (err) {
      console.error("Realtime WS error:", err);
    }
  };

  const handleMessageSession = (event) => {
    try {
      const data = JSON.parse(event.data); // {'id': 'name'}
      const newPlayers = new Map();
      for (const [id, name] of Object.entries(data)) {
        if (name !== "Admin") {
          newPlayers.set(id, name);
        }
      }
      setPlayers(newPlayers);
    } catch (err) {
      console.error("Session WS error:", err);
    }
  };

  return (
    <div className="wait-player-wrapper">
      <div className="player-left-side">
        <h1 className="waiting-title">
          Waiting for your awesome <span className="highlight">crew</span>...
        </h1>
        <div className="players-grid">
          {Array.from(players.entries()).map(([id, name]) => (
            <div key={id} className="player-card">
              {name}
            </div>
          ))}
        </div>
      </div>
      <div className="player-right-side">
        <div className="session-code"># {sessionCode}</div>
        <div className="session-count">👤 {players.size}/40</div>
        <button className="leave-btn" onClick={endSession}>
          ← Leave
        </button>
      </div>
    </div>
  );
};

export default WaitGameStartPlayer;
