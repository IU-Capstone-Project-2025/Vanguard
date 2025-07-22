import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext";
import styles from './styles/WaitGameStartPlayer.module.css';

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

    if (sessionCode) {
      sessionStorage.setItem("sessionCode", sessionCode);
    }

    const handleMessageRealtime = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "next_question") {
          navigate(`/game-process/${sessionCode}`);
        } else if (data.type === "end_session") {
          endSession();
        }
      } catch (err) {
        console.error("Realtime WS error:", err);
      }
    };

    const handleMessageSession = (event) => {
      try {
        const data = JSON.parse(event.data);
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

    const setupWebSockets = () => {
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
    };

    setupWebSockets();

    return () => {
      if (wsRefSession.current) wsRefSession.current.onmessage = null;
      if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
    };
  }, [navigate, sessionCode, connectSession, connectRealtime]);

  const endSession = () => {
    sessionStorage.removeItem("sessionCode");
    sessionStorage.removeItem("nickname");
    closeWsRefRealtime();
    closeWsRefSession();
    navigate("/");
  };

  return (
    <div className={styles['wait-player-wrapper']}>
      <div className={styles['player-left-side']}>
        <h1 className={styles['waiting-title']}>
          Waiting for your awesome <span className={styles['highlight']}>crew</span>...
        </h1>
        <div className={styles['players-grid']}>
          {Array.from(players.entries()).map(([id, name]) => (
            <div key={id} className={styles['player-card']}>
              {name}
            </div>
          ))}
        </div>
      </div>
      <div className={styles['player-right-side']}>
        <div className={styles['session-code']}># {sessionCode}</div>
        <div className={styles['session-count']}>👤 {players.size}/40</div>
        <button className={styles['leave-btn']} onClick={endSession}>
          ← Leave
        </button>
      </div>
    </div>
  );
};

export default WaitGameStartPlayer;