import React, { useState, useEffect } from "react";
import GameProcess from "./GameProcess";

const GameController = ({ sessionId }) => {
  const [question, setQuestion] = useState(null);
  const [options, setOptions] = useState([]);
  const [loading, setLoading] = useState(true);

  // Допустим, что приход данных реализуем REST или WebSocket
  // Здесь делаем простую имитацию REST-запроса
  useEffect(() => {
    async function fetchQuestion() {
      try {
        const res = await fetch(`/api/sessions/${sessionId}/next_question`);
        const data = await res.json();
        setQuestion(data.question);
        setOptions(data.options);
        setLoading(false);
      } catch (error) {
        console.error(error);
      }
    }
    fetchQuestion();
  }, [sessionId]);

  const handleAnswer = async (selectedOption) => {
    // Отправим ответ
    try {
      await fetch(`/api/sessions/${sessionId}/answer`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ answer: selectedOption }),
      });
      // Потом получаем следующий вопрос
      const res = await fetch(`/api/sessions/${sessionId}/next_question`);
      const data = await res.json();
      setQuestion(data.question);
      setOptions(data.options);
    } catch (error) {
      console.error(error);
    }
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
