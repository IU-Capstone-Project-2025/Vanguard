import React, { useState, useEffect } from "react";
import "./styles/QuizStorePage.css";
import sampleImage from "./assets/sampleImage.png"

// Заглушка для получения sessionCode
const handleStartQuiz = () => "123456";

// Функция обработки View
const ViewQuiz = (quizId) => {
  alert(`View quiz with ID: ${quizId}`);
};

const QuizStorePage = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [quizzes, setQuizzes] = useState([]);

  useEffect(() => {
    // Здесь сделай fetch из API, пока просто заглушка
    setQuizzes([
      { id: 1, name: "Frontend Quiz", description: "Test your HTML/CSS skills" },
      { id: 2, name: "Backend Quiz", description: "Server-side challenge" },
      { id: 3, name: "Fullstack Quiz", description: "All-around developer test" },
    ]);
  }, []);

  const filteredQuizzes = quizzes.filter((quiz) =>
    quiz.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="quiz-store-page">
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
                <img src={quiz.imageURL && sampleImage}/>
                <h3>{quiz.name}</h3>
                <p>{quiz.description}</p>
                <div className="quiz-buttons">
                <button className="secondary-button" onClick={() => ViewQuiz(quiz.id)}>View</button>
                <button className="primary-button" onClick={() => {handleStartQuiz()}}>Start</button>
                </div>
            </div>
            ))}
        </div>
        </div>
    </div>
  );
};

export default QuizStorePage;
