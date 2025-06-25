import React, { useState, useEffect } from 'react';
import './styles/GameProcess.css';
import { useParams } from 'react-router-dom';

const GameProcessAdmin = () => {
  const [quiz, setQuiz] = useState(null);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answered, setAnswered] = useState(0);
  const { sessionCode } = useParams();
  const [serverEndpoint, setServerEndpoint] = useState(sessionStorage.getItem('WSEndpoint') || 'http://localhost:8001');
  const [membersNumber, setMembersNumber] = useState(0);
  const [isShowCorrect, setIsShowCorrect] = useState(false);

  console.log("Server Endpoint:", serverEndpoint);
  console.log("Session Code:", sessionCode);
  console.log("Selected Quiz ID:", sessionStorage.getItem('selectedQuizId'));
  console.log(quiz);

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

  const showCorrectAnswer = (correctAnswer) => {
    setIsShowCorrect(true);
    document.querySelectorAll('.option-button').forEach(button => {
      if (button.textContent === correctAnswer) {
        button.classList.add('correct');
      } else {
        button.classList.add('incorrect');
      }
    });
  };

  const hideCorrectAnswer = () => {
    setIsShowCorrect(false);
    document.querySelectorAll('.option-button').forEach(button => {
      button.classList.remove('correct', 'incorrect');
    });
  };

  if (!quiz) {
    return <div className="game-process">Loading Data...</div>;
  }

  const currentQuestion = quiz.questions[currentQuestionIndex];

  return (
    <div className="game-process">
      <h1>{quiz.title}</h1>
      <p>Question {currentQuestionIndex + 1} / {quiz.questions.length}</p>
      <p>Answered: {answered} / {membersNumber}</p>
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
          <button key={index} className={
              option.is_correct && isShowCorrect
               ? "option-button correct" : "option-button"
            }
          >
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
          <button
            onClick={() => {
              showCorrectAnswer(currentQuestion.options.find(option => option.is_correct).text);
            }}
            className="nav-button"
          >
            Show Answer
          </button>
        {currentQuestionIndex < quiz.questions.length - 1 && (
          <button
            onClick={() => {
              setCurrentQuestionIndex(currentQuestionIndex + 1);
              hideCorrectAnswer();
            }}
            className="nav-button"
          >
            Next
          </button>
        )}
      </div>
    </div>
  );
};

export default GameProcessAdmin;
