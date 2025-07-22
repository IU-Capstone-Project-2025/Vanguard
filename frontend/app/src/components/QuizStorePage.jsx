import React, { useState, useEffect } from "react";
import Cookies from "js-cookie";
import { API_ENDPOINTS } from '../constants/api';
import styles from './styles/QuizStorePage.module.css';
import sampleImage from "./assets/sampleImage.png";
import { useNavigate } from "react-router-dom";
import QuizPreviewModal from "./childComponents/QuizPreviewModal";

const QuizStorePage = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [quizzes, setQuizzes] = useState([]);
  const user_id = Cookies.get("user_id");
  const navigate = useNavigate();
  const [selectedQuiz, setSelectedQuiz] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [clickCoordinates, setClickCoordinates] = useState([0, 0]);

  const handleViewQuiz = (e, quizId) => {
    e.preventDefault();
    const quiz = quizzes.find((q) => q.id === quizId);
    setSelectedQuiz(quiz);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedQuiz(null);
  };

  const fetchQuizzes = async () => {
    const url = `${API_ENDPOINTS.QUIZ}/`;
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`Network error: ${response.status}`);
      } 
      const data = await response.json();
      if (!Array.isArray(data)) {
        throw new Error("Expected an array of quizzes.");
      }
      setQuizzes(data);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchQuizzes();
  }, []);

  const createSession = async (quizId) => {
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

    return await response.json();
  };

  const handleStartQuiz = async (quiz) => {
    try {
      const response = await createSession(quiz.id);
      sessionStorage.setItem('selectedQuizId', quiz.id);
      sessionStorage.setItem('sessionCode', response.sessionId);
      sessionStorage.setItem('jwt', response.jwt);
      navigate(`/sessionAdmin/${response.sessionId}`);
    } catch (error) {
      console.error("Error starting quiz:", error);
    }
  };

  const filteredQuizzes = quizzes.filter((quiz) =>
    quiz.title.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleEditQuiz = (quiz) => {
    alert('Feature coming soon');
  };

  return (
    <div className={styles['quiz-store-page']}>
      {isModalOpen && (
        <QuizPreviewModal 
          quiz={selectedQuiz} 
          onClose={closeModal}
          coordinates={clickCoordinates}
          onStart={() => handleStartQuiz(selectedQuiz)}
          onEdit={() => handleEditQuiz(selectedQuiz)}
        />
      )}
      
      <div className={styles['quiz-store-container']}>
        <nav className={styles['nav-links']}>
          <a href="/">Home</a>
          <a href="/join">Join</a>
          <a href="/create">Create</a>
        </nav>

        <div className={styles['top-bar']}>
          <input
            type="text"
            placeholder="Search quizzes..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className={styles['search-input']}
          />
        </div>

        <div className={styles['quiz-list']}>
          {filteredQuizzes.map((quiz) => (
            <div key={quiz.id} className={styles['quiz-card']}>
              <img 
                src={quiz.imageURL ?? sampleImage} 
                alt="Quiz Preview" 
                className={styles['quiz-image']}
              />
              <h3 className={styles['quiz-title']}>{quiz.title}</h3>
              <p className={styles['quiz-description']}>
                {quiz.description === '' ? "No description available" : quiz.description}
              </p>
              <div className={styles['quiz-buttons']}>
                <button
                  className={styles['secondary-button']}
                  onClick={(e) => handleViewQuiz(e, quiz.id)}
                >
                  View
                </button>
                <button
                  className={styles['primary-button']}
                  onClick={() => handleStartQuiz(quiz)}
                >
                  Start
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      <button
        className={styles['floating-create-button']}
        onClick={() => navigate("/constructor/new")}
      >
        ⤬
      </button>
    </div>
  );
};

export default QuizStorePage;