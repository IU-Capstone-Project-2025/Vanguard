/* eslint-disable no-unused-vars */
// QuizStorePage.jsx
import React, { useState, useEffect } from "react";
import Cookies from "js-cookie"
import { API_ENDPOINTS } from '../constants/api';

import "./styles/QuizStorePage.css";
import sampleImage from "./assets/sampleImage.png";
import { useNavigate } from "react-router-dom";
import QuizPreviewModal from "./childComponents/QuizPreviewModal";


const QuizStorePage = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [quizzes, setQuizzes] = useState([]);
  const user_id = Cookies.get("user_id")
  const navigate = useNavigate();
  const [selectedQuiz, setSelectedQuiz] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [clickCoordinates, setClickCoordinates] = useState([0, 0])

  const handleViewQuiz = (e, quizId) => {
    e.preventDefault();
    console.log(e)
    // setClickCoordinates([e.pageX, e.pageY])
    // console.log([e.screenX, e.screenY])
    const quiz = quizzes.find((q) => q.id === quizId);
    setSelectedQuiz(quiz);
    setIsModalOpen(true);
    console.log(clickCoordinates)
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedQuiz(null);
  };

  useEffect(() => {
    const fetchQuizzes = async () => {
      // const queryParams = new URLSearchParams(
      //   {
      //     user_id: user_id,
      //   }
      // )

      const url = `${API_ENDPOINTS.QUIZ}/ `
      // console.log(url)
      try {
        const response = await fetch(url) // заглушка, убрано чтобы не было ошибки
        // console.log(response)
        if (!response.ok) {
          throw new Error(`Network error: ${response.status}`);
          // console.log(data)
        } 
        const data = await response.json();
        if (!Array.isArray(data)) {
          throw new Error("Expected an array of quizzes.");
        }
        setQuizzes(data);
      } catch (error) {
        console.error(error)
      }

    };
    // Заглушка с локальными квизами
    fetchQuizzes();
  }, [quizzes, setQuizzes, user_id]);

  const createSession = async (quizId) => {
      console.log("Creating session with quizId:", quizId, "userName:", Cookies.get("user_nickname"));
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
      return data; // возвращает объект вида: {"serverWsEndpoint": "string","jwt": "string", "sessionId":"string"}
    };

  const handleStartQuiz = async (quiz) => {
    const response = await createSession(quiz.id);

    sessionStorage.setItem('selectedQuizId', quiz.id);
    sessionStorage.setItem('sessionCode', response.sessionId);
    sessionStorage.setItem('jwt', response.jwt);

    navigate(`/ask-to-join/${response.sessionId}`);
  };

  // Локальный поиск по имени квиза
  const filteredQuizzes = quizzes.filter((quiz) =>
    quiz.title.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleEditQuiz = (quiz) => {
    alert('unable yet')
    // navigate(`/edit-quiz/${quiz.id}`)
  }


  return (
    <div className="quiz-store-page">
      {isModalOpen && (
        <QuizPreviewModal 
          quiz={selectedQuiz} 
          onClose={closeModal} 
          coordinates={clickCoordinates}
          onStart={() => handleStartQuiz(selectedQuiz)}
          onEdit={() => handleEditQuiz(selectedQuiz)}
        />
      )}
      <div className="quiz-store-container">
        <nav className="nav-links">
          <a href="/">Home</a>
          <a href="/join">Join</a>
          <a href="/create">Create</a>
        </nav>

        <div className="top-bar">
          <input
            type="text"
            placeholder="Search quizzes..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div className="quiz-list">
          {filteredQuizzes.map((quiz) => (
            <div key={quiz.id} className="quiz-card">
              <img src={quiz.imageURL ?? sampleImage} alt="Quiz Preview" />
              <h3>{quiz.title}</h3>
              <p>{quiz.description}</p>
              <div className="quiz-buttons">
                <button
                  className="secondary-button"
                  onClick={(e) => handleViewQuiz(e, quiz.id)}
                >
                  View
                </button>
                <button
                  className="primary-button"
                  onClick={(e) => handleStartQuiz(quiz)}
                >
                  Start
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default QuizStorePage;
