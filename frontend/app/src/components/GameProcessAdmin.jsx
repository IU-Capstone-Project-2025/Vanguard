import React, { useEffect, useState } from 'react';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import './styles/GameProcess.css'
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../constants/api';
import ShapedButton from './childComponents/ShapedButton';
import ShowLeaderBoardComponent from './childComponents/ShowLeaderBoardComponent';
import Alien from './assets/Alien.svg'
import Corona from './assets/Corona.svg'
import Ghosty from './assets/Ghosty.svg'
import Cookie6 from './assets/Cookie6.svg'

const GameProcessAdmin = () => {
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const [currentQuestion, setCurrentQuestion] = useState(sessionStorage.getItem('currentQuestion') != undefined ?
    JSON.parse(sessionStorage.getItem('currentQuestion')) : {});
  const [questionIndex, setQuestionIndex] = useState(0);
  const [leaderboardVisible, setLeaderboardVisible] = useState(false);
  const [leaderboardData, setLeaderboardData] = useState(null);

  const questionsAmount = currentQuestion.questionsAmount - 1;
  const navigate = useNavigate();
  const [questionOptions] = useState([
    Alien,
    Corona,
    Ghosty,
    Cookie6
  ]);

  useEffect(() => {
    const token = sessionStorage.getItem('jwt');
    if (!token) return;

    return () => {
      if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
      if (wsRefSession.current) wsRefSession.current.onmessage = null;
      if (wsRefRealtime.current) wsRefRealtime.current.onclose = { finishSession };
      if (wsRefSession.current) wsRefSession.current.onclose = { finishSession };
    };
  }, [connectSession, connectRealtime, wsRefSession, wsRefRealtime]);

  const toNextQuestion = async (sessionCode) => {
    if (!sessionCode) {
      console.error('Session code is not available');
      return;
    }
    try {
      const response = await fetch(`${API_ENDPOINTS.SESSION}/session/${sessionCode}/nextQuestion`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) throw new Error('Failed to start next question');
    } catch (error) {
      console.error('Error starting next question:', error);
    }
  };

  const listenQuizQuestion = async (sessionCode) => {
    if (!sessionCode) {
      console.error('Session code is not available');
      return;
    }

    wsRefRealtime.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'question') {
        setCurrentQuestion(data);
        sessionStorage.setItem('currentQuestion', JSON.stringify(data));
      } else if (data.type === 'leaderboard') {
        console.log('Received leaderboard:', data);
        sessionStorage.setItem('leaders', JSON.stringify(data.payload.users.slice(0, 3)));
        setLeaderboardData(data.payload);
        setLeaderboardVisible(true);
      }
    };
  };

  const finishSession = async (code) => {
    await toNextQuestion(code);
    await listenQuizQuestion(code);
    // sessionStorage.removeItem('sessionCode');
    sessionStorage.removeItem('quizData');
    sessionStorage.removeItem('currentQuestion');

    closeWsRefRealtime();
    closeWsRefSession();
    navigate('/final');
    // если leaderboard не придёт, не вызываем завершение сразу
  };

  const handleLeaderboardClick = () => {
    wsRefRealtime.current.send(JSON.stringify({ type: 'next_question' }));
    setLeaderboardVisible(false);
  }

  const handleNextQuestion = async (e) => {
    e.preventDefault();
    const sessionCode = sessionStorage.getItem('sessionCode');

    await toNextQuestion(sessionCode);
    setQuestionIndex((prevIndex) => prevIndex + 1);
    await listenQuizQuestion(sessionCode);
  };

  return (
    <div className="game-process">
      {leaderboardVisible && leaderboardData ? (
        <ShowLeaderBoardComponent
          leaderboardData={leaderboardData}
          onClose={() => handleLeaderboardClick()} 
        />
      ) : (
        <>
          <div className="controller-question-title">
            <img src={currentQuestion.payload} alt="Question" className="question-image" height={300} />
            <h2>{currentQuestion ? currentQuestion.text : 'Waiting for question…'}</h2>
          </div>

          <div className="options-grid">
            {currentQuestion && currentQuestion.options.map((option, idx) => (
              <ShapedButton
                key={idx}
                shape={questionOptions[idx]}
                text={option.text}
                onClick={() => console.log('svg clicked')}
              />
            ))}
          </div>

          <div className="button-group">
            {questionIndex < questionsAmount && (
              <button onClick={handleNextQuestion} className="button">
                Next
              </button>
            )}
            <button onClick={() => finishSession(sessionStorage.getItem('sessionCode'))} className="nav-button">
              Finish
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default GameProcessAdmin;
