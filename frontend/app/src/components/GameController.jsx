import React, { useState, useEffect } from "react";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext";
import { useNavigate } from "react-router-dom";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import ShapedButton from "./childComponents/ShapedButton";
import ShowQuizStatistics from "./childComponents/ShowQuizStatistics";
import "./styles/GameProcess.css";
import Alien from './assets/Alien.svg';
import Corona from './assets/Corona.svg';
import Ghosty from './assets/Ghosty.svg';
import Cookie6 from './assets/Cookie6.svg';

const GameController = () => {
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const { wsRefSession, closeWsRefSession } = useSessionSocket();
  const [popularAnswers, setPopularAnswers] = useState({});
  const [userAnswers, setUserAnswers] = useState([]);
  const [question, setQuestion] = useState({
    options: [Alien, Corona, Ghosty, Cookie6]
  });
  const [hasAnswered, setHasAnswered] = useState(false);
  const [showStatistics, setShowStatistics] = useState(false);
  const sessionCode = sessionStorage.getItem("sessionCode");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const endSession = React.useCallback(() => {
    console.log(`Ending session... ${sessionCode}`);
    sessionStorage.removeItem('sessionCode');
    sessionStorage.removeItem('nickname');
    closeWsRefRealtime();
    closeWsRefSession();
    navigate('/');
  }, [sessionCode, closeWsRefRealtime, closeWsRefSession, navigate]);

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
          setShowStatistics(false);
        }

        if (data.type === "question_stat") {
          console.log("Received question statistics:", data);
          setPopularAnswers(data.payload.answers);
          setShowStatistics(true);
          setLoading(false);

          // Обновляем массив правильности ответов
          setUserAnswers((prev) => {
            const updated = [...prev, !!data.correct];
            sessionStorage.setItem("userAnswers", JSON.stringify(updated));
            return updated;
          });
        }
      };

      wsRefRealtime.current.onclose = () => {
        endSession();
      };
    }

    if (wsRefSession.current) {
      wsRefSession.current.onclose = () => {
        endSession();
      };
    }

    return () => {
      if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
      if (wsRefSession.current) wsRefSession.current.onmessage = null;
    };
  }, [connectRealtime, endSession, wsRefRealtime, wsRefSession]);

  const handleAnswer = (index) => {
    if (!wsRefRealtime.current) return;
    const timestamp = new Date().toISOString();
    console.log(`Sending answer: ${index} at ${timestamp}`);
    wsRefRealtime.current.send(JSON.stringify({ option: index, timestamp }));
    setHasAnswered(true);
  };

  if (loading) {
    return (
      <div style={{ color: "#F9F3EB", padding: "2vw" }}>
        Question loading...
      </div>
    );
  }

  return (
    <div className="game-process-player">
      {showStatistics && (
        <ShowQuizStatistics
          stats={popularAnswers}
          onClose={() => setShowStatistics(false)}
        />
      )}

      {!showStatistics && (
        hasAnswered ? (
          <p className="waiting-text">Waiting for next question...</p>
        ) : (
          <div className="options-grid-player">
            {question.options.map((option, idx) => (
              <ShapedButton
                key={idx}
                shape={option}
                text=""
                onClick={() => handleAnswer(idx)}
              />
            ))}
          </div>
        )
      )}
    </div>
  );
};

export default GameController;
