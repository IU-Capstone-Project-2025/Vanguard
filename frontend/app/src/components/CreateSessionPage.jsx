import React, { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie";
import { API_ENDPOINTS } from '../constants/api';

import styles from './styles/CreateSessionPage.module.css';

const CreateSessionPage = () => {
  const navigate = useNavigate();

  const [selectedQuiz, setSelectedQuiz] = useState(null);
  const [quizzes, setQuizzes] = useState([]);
  const [search, setSearch] = useState("");
  const [sessionCode, setSessionCode] = useState(null);
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const inputRef = useRef(null);

  useEffect(() => {
    const fetchQuizzes = async () => {
      try {
        setIsLoading(true);
        const response = await fetch(`${API_ENDPOINTS.QUIZ}/`);
        if (!response.ok) {
          throw new Error(`Network error: ${response.status}`);
        }
        const data = await response.json();
        if (!Array.isArray(data)) {
          throw new Error("Expected an array of quizzes.");
        }
        setQuizzes(data.map(quiz => ({ id: quiz.id, title: quiz.title })));
      } catch (error) {
        // console.error("Error fetching quizzes:", error);
        setError("Failed to load quizzes. Please try again later.");
      } finally {
        setIsLoading(false);
      }
    };
    fetchQuizzes();
  }, []);

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  const handleQuizSelection = (quiz) => {
    sessionStorage.setItem('selectedQuizId', quiz.id);
    setSelectedQuiz(quiz);
    setError(null);
  };

  const createSession = async (quizId) => {
    try {
      setIsLoading(true);
      const response = await fetch(`${API_ENDPOINTS.SESSION}/sessions`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          "userName": "Admin",
          "quizId": quizId,
        }),
      });

      if (!response.ok) throw new Error("Failed to create session");

      const data = await response.json();
      return data;
    } catch (error) {
      // console.error("Error creating session:", error);
      setError("Failed to create session. Please try again.");
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const handlePlay = async () => {
    if (!selectedQuiz) {
      setError("Please select a quiz first");
      return;
    }

    try {
      const sessionData = await createSession(selectedQuiz.id);
      setSessionCode(sessionData.sessionId);

      sessionStorage.setItem('selectedQuizId', selectedQuiz.id);
      sessionStorage.setItem('sessionCode', sessionData.sessionId);
      sessionStorage.setItem('jwt', sessionData.jwt);

      navigate(`/sessionAdmin/${sessionData.sessionId}`);
    } catch (error) {
      // Error is already handled in createSession
    }
  };

  return (
    <div className={styles['create-session-main-content']}>
      <div className={styles['left-side']}>
        <div className={styles.title}>
          <h2>Ready to roll? <br /> Pick your <span>quiz</span></h2>
          <div className={styles['button-group']}>
            <button
              className={styles['play-button']}
              onClick={handlePlay}
              disabled={!selectedQuiz || isLoading}
            >
              {isLoading ? "Loading..." : "▶ Play"}
            </button>
            <button
              className={styles['enter-store-button']}
              onClick={() => navigate("/store")}
            >
              + Quiz Store
            </button>
          </div>
          {error && <div className={styles.error}>{error}</div>}
        </div>
      </div>

      <div className={styles['right-side']}>
        <div className={styles['quiz-list-container']}>
          <div className={styles['quiz-search-panel']}>
            <input
              type="text"
              ref={inputRef}
              placeholder="Search the quiz"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className={styles['quiz-search-input']}
            />

            <div className={styles['session-quiz-list']}>
              {isLoading ? (
                <div className={styles.loading}>Loading quizzes...</div>
              ) : (
                quizzes
                  .filter((quiz) =>
                    quiz.title.toLowerCase().includes(search.toLowerCase())
                  )
                  .map((quiz) => (
                    <div
                      key={quiz.id}
                      className={`${styles['quiz-item']} ${
                        selectedQuiz?.id === quiz.id ? styles.selected : ''
                      }`}
                      onClick={() => handleQuizSelection(quiz)}
                    >
                      {quiz.title}
                    </div>
                  ))
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateSessionPage;