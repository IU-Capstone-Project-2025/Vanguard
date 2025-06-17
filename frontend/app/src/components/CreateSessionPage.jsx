import React, { useState } from "react";
import './styles/styles.css';
import { useNavigate } from "react-router-dom";

const CreateSessionPage = () => {
    const navigate = useNavigate();
    const [selectedQuiz, setSelectedQuiz] = useState(null);
    const quizzes = ["Quiz 1", "Quiz 2"]; // Placeholder, should be dynamic in real app
    const [search, setSearch] = useState("");

    const handlePlay = () => {
        if (selectedQuiz) {
            // Logic to start the session with selectedQuiz
            navigate('/play'); // You may pass quiz info as state or param
        }
    };

    return (
        <div className="create-session-main-content">
            <div className="left-side">
                <div className="title">
                    <h2>
                        Now choose the quiz <br /> to start a game
                    </h2>
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
                            onClick={() => navigate('/store')}
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
