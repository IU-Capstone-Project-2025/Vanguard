import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

import './styles/styles.css';

const CreateSessionPage = () => {
  const navigate = useNavigate();

  const [selectedQuiz, setSelectedQuiz] = useState(null);
  const [quizzes, setQuizzes] = useState([]);
  const [search, setSearch] = useState("");
  const [sessionCode, setSessionCode] = useState(null);

  useEffect(() => {
    const fetchQuizzes = async () => {
      try {
        const response = await fetch("/api/quiz/");
        if (!response.ok) {
          throw new Error(`Network error: ${response.status}`);
        }
        const data = await response.json();
        if (!Array.isArray(data)) {
          throw new Error("Expected an array of quizzes.");
        }
        // Мапим массив, берем только id и title
        setQuizzes(data.map(quiz => ({ id: quiz.id, title: quiz.title })));
      } catch (error) {
        console.error("Error fetching quizzes:", error);
      }
    };
    fetchQuizzes();
  }, []);

  const handleQuizSelection = (quiz) => {
    sessionStorage.setItem('selectedQuizId', quiz.id);
    setSelectedQuiz(quiz);
  };

  const handlePlay = () => {
    if (selectedQuiz) {
      sessionStorage.setItem('selectedQuizId', selectedQuiz.id);
      sessionStorage.setItem('sessionCode', sessionCode);
      navigate(`/ask-to-join/${sessionCode}`);
    }
  };

  return (
    <div className="create-session-main-content">
      <div className="left-side">
        <div className="title">
          <h2>Now choose the quiz <br /> to start a game</h2>
          <div className="button-group">
            <button
              className="play-button"
              onClick={handlePlay}
              disabled={!selectedQuiz}
            >
              ▶ Play
            </button>
            <button
              className="enter-store-button"
              onClick={() => navigate("/store")}
            >
              + Enter quiz Store
            </button>
          </div>
        </div>
      </div>

      <div className="right-side">
        <div className="quiz-list-container">
          <div className="quiz-search-panel">
            <input
              type="text"
              placeholder="Search the quiz"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="quiz-search-input"
            />

            <div className="quiz-list">
              {quizzes
                .filter((quiz) =>
                  quiz.title.toLowerCase().includes(search.toLowerCase())
                )
                .map((quiz) => (
                  <div
                    key={quiz.id}
                    className={`quiz-item ${selectedQuiz?.id === quiz.id ? 'selected' : ''}`}
                    onClick={() => handleQuizSelection(quiz)}
                  >
                    {quiz.title}
                  </div>
                ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateSessionPage;
