import React, {useState, useEffect} from 'react';
import './styles/GameProcess.css';

const GameProcessAdmin = () => {
    const [quiz, setQuiz] = useState(null);
    const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);

    useEffect(() => {
        // Здесь ты делаешь fetch из бекенда
        const fetchData = async () => {
            let id = sessionStorage.getItem('selectedQuizId'); // Получаем ID квиза из sessionStorage
            const response = await fetch(`/api/quiz/${id}`); // Заменить {id} на реальный ID квиза
            const data = await response.json();
            setQuiz(data);
        };
        fetchData();
    }, []);

    if (!quiz) {
        return <div className="game-process">Loading Data...</div>;
    }
    const toNextQuestion = async (sessionCode) => {

        const response = await fetch(`/api/session/session/${sessionCode}/nextQuestion`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ "code": sessionCode}),
        });

        if (!response.ok) throw new Error("Failed to get to the next question session");

    };
    const currentQuestion = quiz.questions[currentQuestionIndex];
    const handleNextQuestion = async () => {
        setCurrentQuestionIndex(currentQuestionIndex + 1)
        await toNextQuestion(sessionStorage.getItem('selectedQuizId'));
    }
    return (
        <div className="game-process">
            <h1>{quiz.title}</h1>
            <p>Question{currentQuestionIndex + 1} / {quiz.questions.length}</p>
            <div className="question-block">
                <h2>{currentQuestion.text}</h2>
                {currentQuestion.image_url && (
                    <img
                        src={currentQuestion.image_url}
                        alt="Question"
                        className="question-image"
                    />
                )}
            </div>
            <div className="options-grid">
                {currentQuestion.options.map((option, index) => (
                    <button key={index} className="option-button">
                        {option.text}
                    </button>
                ))}
            </div>
            <div className="navigation-buttons">
                {currentQuestionIndex > 0 && (
                    <button
                        onClick={() =>
                            setCurrentQuestionIndex(currentQuestionIndex - 1)
                        }
                        className="nav-button"
                    >
                        Back
                    </button>
                )}
                {currentQuestionIndex < quiz.questions.length - 1 && (
                    <button
                        onClick={handleNextQuestion}
                        className="nav-button">
                        Next
                    </button>
                )}
            </div>
        </div>
    );
};

export default GameProcessAdmin;
