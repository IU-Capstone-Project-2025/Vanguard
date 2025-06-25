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
        const response = await fetch("api/quiz/", {
          method: "GET",
          headers: {
            "Content-Type": "application/json"
          }
        });
        console.log("Response from /api/quiz:", response);
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

    let userId = sessionStorage.getItem('userId');
    if (!userId) {
      userId = crypto.randomUUID();
      sessionStorage.setItem('userId', userId);
    }

    sessionStorage.setItem('selectedQuizId', selectedQuiz.id);
    if (selectedQuiz) {
      fetch(`http://localhost:8081/sessions/`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ quizId: sessionStorage.getItem('selectedQuizId')})
      })
        .then(res => {
          if (!res.ok) throw new Error("Failed to create session");
          return res.json();
        })
        .then(data => {
          if (!data.sessionCode) throw new Error("No sessionCode returned");
          setSessionCode(data.sessionCode);
          sessionStorage.setItem('sessionCode', data.sessionCode);
          navigate(`/ask-to-join/${data.sessionCode}`);
        })
        .catch(err => {
          console.error("Error creating session:", err);
          alert("Failed to create session. Please try again.");
        });
      return;
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
