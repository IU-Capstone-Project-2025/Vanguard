import React, { useEffect, useState } from "react";
import './styles/styles.css';
import { useNavigate } from "react-router-dom";

const CreateSessionPage = () => {
    const navigate = useNavigate();
    const [selectedQuiz, setSelectedQuiz] = useState(null);
    const [quizzes, setQuizzes] = useState([]);
    const [search, setSearch] = useState("");
    const [sessionCode, setSessionCode] = useState(null);

    useEffect(() => {
        const fetchQuizzes = async () => {
            try {
                const response = await fetch("/api/quiz/"); // <-- относительный путь
                if (!response.ok) {
                    throw new Error(`Network error: ${response.status}`);
                }
                const data = await response.json();
                if (!Array.isArray(data)) {
                    throw new Error("Expected an array of quizzes.");
                }
                setQuizzes(data.map(quiz => quiz.title)); // <-- Здесь берем title
            } catch (error) {
                console.error("Error fetching quizzes:", error);
            }
        };
        fetchQuizzes();
    }, []);

    const handlePlay = () => {
        if (selectedQuiz) {
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
                                .filter((quiz) => quiz.toLowerCase().includes(search.toLowerCase()))
                                .map((quiz, index) => (
                                    <div
                                        key={index}
                                        className={`quiz-item ${selectedQuiz === quiz ? 'selected' : ''}`}
                                        onClick={() => setSelectedQuiz(quiz)}
                                    >
                                        {quiz}
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
