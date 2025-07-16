import React, { useState, useEffect, useCallback } from "react";
import { useRealtimeSocket } from "../contexts/RealtimeWebSocketContext";
import { useSessionSocket } from "../contexts/SessionWebSocketContext";
import { useNavigate } from "react-router-dom";
import ShapedButton from "./childComponents/ShapedButton";
import ShowQuizStatistics from "./childComponents/ShowQuizStatistics";
import "./styles/GameProcess.css";
import Alien from "./assets/Alien.svg";
import Corona from "./assets/Corona.svg";
import Ghosty from "./assets/Ghosty.svg";
import Cookie6 from "./assets/Cookie6.svg";

const GameController = () => {
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const { wsRefSession, closeWsRefSession } = useSessionSocket();
  const navigate = useNavigate();

  const [question, setQuestion] = useState({
    options: [Alien, Corona, Ghosty, Cookie6]
  });

  const [choosenOption, setChosenOption] = useState(null);
  const [popularAnswers, setPopularAnswers] = useState({});
  const [userAnswers, setUserAnswers] = useState([]);
  const [stage, setStage] = useState("waiting_question"); // "question", "waiting", "statistics", "waiting_question"

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

    if (!token || !sessionCode) return;

    connectRealtime(token, sessionCode);

    if (wsRefRealtime.current) {
      wsRefRealtime.current.onmessage = (event) => {
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
            break;

          case "question_stat":
            console.log("Received question statistics:", data);
            setPopularAnswers(data.payload.answers);
            setStage("statistics");

            // Обновляем массив правильности ответов
            setUserAnswers((prev) => {
              const updated = [...prev, !!data.correct];
              sessionStorage.setItem("userAnswers", JSON.stringify(updated));
              return updated;
            });
            break;

          default:
            break;
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
    setChosenOption(index);
    if (!wsRefRealtime.current) return;

    const timestamp = new Date().toISOString();
    const answerMessage = {
      type: "answer",
      option: index,
      timestamp
    };

    wsRefRealtime.current.send(JSON.stringify(answerMessage));
    setStage("waiting");
  };

  return (
    <div className="game-process-player">
      {stage === "statistics" && (
        <ShowQuizStatistics
          stats={popularAnswers}
          onClose={() => setStage("waiting_question")} // переходим в ожидание next_question
        />
      )}

      {stage === "question" && (
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
      )}

      {stage === "waiting" && (
        <div className="options-grid-player">
          <ShapedButton
            key={0}
            shape={question.options[choosenOption]}
            text=" "
            onClick={() => {}}
          />
        </div>
      )}

      {stage === "waiting_question" && (
        <div className="waiting-message">Waiting for the next question...</div>
      )}
    </div>
  );
};

export default GameController;
