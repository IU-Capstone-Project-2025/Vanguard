import React, { useState, useEffect } from "react";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext";
import "./styles/GameProcess.css";
import { useNavigate } from "react-router-dom";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";

const GameController = () => {
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const {wsRefSession, closeWsRefSession} = useSessionSocket();
  const [question, setQuestion] = useState(
    {"options": [
      "⬣",
      "⬥",
      "✠",
      "❇"
    ]} // Default empty question to avoid errors
  )
  const [hasAnswered, setHasAnswered] = useState(false);
  const sessionCode = sessionStorage.getItem("sessionCode");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = sessionStorage.getItem("jwt");
    const sessionCode = sessionStorage.getItem("sessionCode");

    if (!token || !sessionCode) return;

    connectRealtime(token, sessionCode);

    if (wsRefRealtime.current) {
      wsRefRealtime.current.onmessage = (event) => {
        const data = JSON.parse(event.data);
        
        if (data.type === "end") {
          endSession();
        }
        if (data.type === "next_question") {
          setHasAnswered(false);
          setLoading(false);
        }
        
      };

      wsRefRealtime.current.onclose = () => {
        endSession();
      }
      wsRefSession.current.onclose = () => {
        endSession();
      }
    }

    return () => {
      if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
    };
  }, [connectRealtime, wsRefRealtime]);

  const endSession = () => {
    console.log(`Ending session... ${sessionCode}`);
    sessionStorage.removeItem('sessionCode');
    closeWsRefRealtime();
    closeWsRefSession();
    navigate('/');
  }

  const handleAnswer = (index) => {
    if (!wsRefRealtime.current) return;
    wsRefRealtime.current.send(JSON.stringify({ option: index }));
    setHasAnswered(true);
  };

  if (loading) {
    return (
      <div style={{ color: "#F9F3EB", padding: "2vw" }}>
        Загрузка вопроса...
      </div>
    );
  }

  return (
    <div className="game-process-container">
      {hasAnswered ? (
        <p className="waiting-text">Ожидание следующего вопроса…</p>
      ) : (
        <div className="options-grid">
          {question.options.map((option, idx) => (
            <button
              key={idx}
              className="option-button"
              onClick={() => handleAnswer(idx)}
            >
              {option}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default GameController;
