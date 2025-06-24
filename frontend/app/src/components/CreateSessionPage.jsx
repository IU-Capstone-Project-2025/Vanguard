import React, { useEffect, useState } from "react";
import './styles/styles.css';
import { useNavigate } from "react-router-dom";

const CreateSessionPage = () => {
    const navigate = useNavigate();
    const [selectedQuiz, setSelectedQuiz] = useState(null);
    const [quizzes, setQuizzes] = useState([])
    const [search, setSearch] = useState("");
    const [SessionCode, setSessionCode] = useState(); // Mocked session code, should be generated dynamically

    useEffect(() => {
        // Mock fetching quizzes from an API or database
        const url = 'http://localhost:8000/api/quiz/'; // Replace with your actual API endpoint
        const fetchQuizzes = async () => {
            try {
                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data = await response.json();
                setQuizzes(data.map(quiz => quiz.name));
                console.log(quizzes) // Assuming the API returns an array of quiz objects with a 'name' property
            } catch (error) {
                console.error('Error fetching quizzes:', error);
            }
        };
        fetchQuizzes();}
    )

    const handlePlay = () => {
        if (selectedQuiz) {
            // Logic to start the session with selectedQuiz

            navigate(`/ask-to-join/${SessionCode}`); // You may pass quiz info as state or param
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
