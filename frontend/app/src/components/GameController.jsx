import React, { useState, useEffect } from "react";
import GameProcess from "./GameProcessPage";

const GameController = ({ sessionId }) => {
  const [question, setQuestion] = useState(null);
  const [options, setOptions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    // Подключаем WebSocket
    const ws = new WebSocket(`wss://example.com/sessions/${sessionId}`);
    setSocket(ws);

    ws.onopen = () => {
      console.log("✅ WebSocket connected");
      // Дополнительно можно отправить сообщение для инициализации
      ws.send(JSON.stringify({ action: "join", sessionId }));
    };
    ws.onerror = (error) => {
      console.error("❌ WebSocket error:", error);
    };
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.type === "new_question") {
        setQuestion(data.question);
        setOptions(data.options);
        setLoading(false);
      } else if (data.type === "end") {
        setQuestion(null);
        setOptions([]);
      }
    };
    ws.onclose = () => {
      console.log("ℹ️ WebSocket closed");
    };
    return () => {
      ws.close();
    };
  }, [sessionId]);

  const handleAnswer = (selectedOption) => {
    if (!socket) return;

    const message = {
      action: "submit_answer",
      sessionId,
      answer: selectedOption,
    };
    socket.send(JSON.stringify(message));
  };

  if (loading) {
    return (
      <div style={{ color: "#F9F3EB", padding: "2vw" }}>
        Загрузка вопроса...
      </div>
    );
  }

  if (!question) {
    return (
      <div style={{ color: "#F9F3EB", padding: "2vw" }}>
        Игра окончена или нет активных вопросов.
      </div>
    );
  }

  return (
    <GameProcess
      question={question}
      options={options}
      onAnswer={handleAnswer}
    />
  );
};

export default GameController;
