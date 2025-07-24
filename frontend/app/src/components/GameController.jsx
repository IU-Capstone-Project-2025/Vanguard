import React, { useState, useEffect, useCallback } from "react";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import { useNavigate } from "react-router-dom";
import ShapedButton from "./childComponents/ShapedButton";
import ShowQuizStatistics from "./childComponents/ShowQuizStatistics";
import styles from './styles/GameController.module.css';
import PentagonYellow from './assets/Pentagon-yellow.svg';
import CoronaIndigo from './assets/Corona-indigo.svg';
import ArrowOrange from './assets/Arrow-orange.svg';
import Cookie4Blue from './assets/Cookie4-blue.svg';

const GameController = () => {
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const { wsRefSession, closeWsRefSession } = useSessionSocket();
  const navigate = useNavigate();
  const [options, setOptions] = useState([])

  const [question, setQuestion] = useState({
    options: [PentagonYellow, CoronaIndigo, ArrowOrange, Cookie4Blue]
  });

  const [choosenOption, setChosenOption] = useState(null);
  const [correct, setCorrect] = useState(false);
  const [popularAnswers, setPopularAnswers] = useState({});
  const [userAnswers, setUserAnswers] = useState([]);
  const [stage, setStage] = useState("question"); // "question", "waiting", "statistics", "waiting_question"
  const [error, setError] = useState(null);

  const sessionCode = sessionStorage.getItem("sessionCode");

  const endSession = useCallback(() => {
    console.log(`Ending session... ${sessionCode}`);
    sessionStorage.removeItem("sessionCode");
    sessionStorage.removeItem("nickname");
    closeWsRefRealtime();
    closeWsRefSession();
    navigate("/");
  }, [sessionCode, closeWsRefRealtime, closeWsRefSession, navigate]);

  useEffect(() => {
    const token = sessionStorage.getItem("jwt");
    const sessionCode = sessionStorage.getItem("sessionCode");

    if (!token || !sessionCode) {
      setError("Missing session credentials");
      return;
    }

    try {
      connectRealtime(token, sessionCode);
    } catch (err) {
      setError("Failed to connect to the game server");
      console.error("Connection error:", err);
    }

    const handleRealtimeMessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("Realtime message received:", data);

        switch (data.type) {
          case "end":
            endSession();
            break;

          case "next_question":
            setChosenOption(null);
            setPopularAnswers({});
            setStage("question");
            setCorrect(false);
            break;
            
            case "question_stat":
              setPopularAnswers(data.payload.answers);
              setCorrect(data.correct);
              setOptions(data.options)  
              setStage("statistics");
              
              setUserAnswers((prev) => {
              const updated = [...prev, !!data.correct];
              sessionStorage.setItem("userAnswers", JSON.stringify(updated));
              return updated;
            });
            break;

          case "next_message":
            setStage("question");
            setChosenOption(null);
            break;

          default:
            break;
        }
      } catch (err) {
        console.error("Error processing message:", err);
        setError("Error processing game data");
      }
    };

    if (wsRefRealtime.current) {
      wsRefRealtime.current.onmessage = handleRealtimeMessage;
      wsRefRealtime.current.onclose = endSession;
    }

    return () => {
      if (wsRefRealtime.current) {
        wsRefRealtime.current.onmessage = null;
        wsRefRealtime.current.onclose = null;
      }
    };
  }, [connectRealtime, endSession, wsRefRealtime]);

  const handleAnswer = (index) => {
    setChosenOption(index);
    if (!wsRefRealtime.current) {
      setError("Connection to game server lost");
      return;
    }

    try {
      const timestamp = new Date().toISOString();
      const answerMessage = {
        type: "answer",
        option: index,
        timestamp
      };
      wsRefRealtime.current.send(JSON.stringify(answerMessage));
      setStage("waiting");
    } catch (err) {
      console.error("Error sending answer:", err);
      setError("Failed to submit answer");
    }
  };

  return (
    <div className={styles['game-process-player']}>
      {error && <div className={styles.error}>{error}</div>}

      {stage === "statistics" && (
        <ShowQuizStatistics
          stats={popularAnswers}
          correct={correct}
          options={options}
          onClose={() => setStage("question")}
        />

      )}

      {stage === "question" && (
        <div className={styles['options-grid-player']}>
          {question.options.map((option, idx) => (
            <div 
              key={idx} 
              className={`${styles['controller-answer-option']} ${
                idx === 0 || idx === 2 ? styles.left : styles.right
              }`}
            >
              <button
                className={styles['option-button']}
                onClick={() => handleAnswer(idx)}
              >
                <img src={question.options[idx]} alt={`option ${idx + 1}`} />
              </button>
            </div>
          ))}
        </div>
      )}

      {stage === "waiting" && (
        <div className={styles['options-grid-player-waiting']}>
          <img src={question.options[choosenOption]} alt="selected shape" />
        </div>
      )}

      {stage === "waiting_question" && (
        <div className={styles['waiting-message']}>
          Waiting for the next question...
        </div>
      )}
    </div>
  );
};

export default GameController;