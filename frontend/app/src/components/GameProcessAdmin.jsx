import React, { useEffect, useState } from 'react';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import './styles/GameProcess.css';
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../constants/api';
import ShapedButton from './childComponents/ShapedButton';
import ShowLeaderBoardComponent from './childComponents/ShowLeaderBoardComponent';
import PentagonYellow from './assets/Pentagon-yellow.svg';
import CoronaIndigo from './assets/Corona-indigo.svg';
import ArrowOrange from './assets/Arrow-orange.svg';
import Cookie4Blue from './assets/Cookie4-blue.svg';
import Triangle from './assets/Triangle.svg';
import Alien from './assets/Alien.svg';

const GameProcessAdmin = () => {
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const [currentQuestion, setCurrentQuestion] = useState(sessionStorage.getItem('currentQuestion') != undefined ?
    JSON.parse(sessionStorage.getItem('currentQuestion')) : {});
  const [questionIndex, setQuestionIndex] = useState(1);
  const [leaderboardVisible, setLeaderboardVisible] = useState(false);
  const [leaderboardData, setLeaderboardData] = useState(null);

  const questionsAmount = currentQuestion.questionsAmount;
  const navigate = useNavigate();
  const [questionOptions] = useState([
    PentagonYellow,
    CoronaIndigo,
    ArrowOrange,
    Cookie4Blue
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
        headers: { 'Content-Type': 'application/json' },
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
      if (data.type === 'leaderboard') {
        sessionStorage.setItem('leaders', JSON.stringify(data.payload.users));
        setLeaderboardData(data.payload);
        setLeaderboardVisible(true);
      } else if (data.type === 'question') {
        setCurrentQuestion(data);
        setQuestionIndex(data.questionId);
        sessionStorage.setItem('currentQuestion', JSON.stringify(data));
      }
    };
  };

  const finishSession = async (code) => {
    sessionStorage.removeItem('quizData');
    sessionStorage.removeItem('currentQuestion');
    if (!wsRefRealtime || !wsRefSession) {
      navigate('/');
    } else {
      await toNextQuestion(code);
      await listenQuizQuestion(code);
      closeWsRefRealtime();
      closeWsRefSession();
      navigate('/final');
    }
  };

  const handleLeaderboardClick = () => {
    wsRefRealtime.current.send(JSON.stringify({ type: 'next_question' }));
    setLeaderboardVisible(false);
  };

  const handleNextQuestion = async (e) => {
    e.preventDefault();
    const sessionCode = sessionStorage.getItem('sessionCode');
    
    await toNextQuestion(sessionCode);
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
        <div className="game-container">
          <div className="question-section">
            <div className="question-header">
              <div className='answer-indicator'>
                <span className='indicator-text'>11/17</span>
              </div>

              {currentQuestion.payload && (
                <img
                  src={currentQuestion.payload}
                  alt="Question"
                  className="process-question-image"
                />
              )}

              <div className='timer-indicator'>
                <span className='indicator-text'>11/17</span>
              </div>
            </div>
            <div className="question-body">
              <div className="question-number-bubble">
                {questionIndex}
              </div>
              <h2 className="question-text">
                {currentQuestion?.text || 'Waiting for question…'}
              </h2>
              
              <div className="question-go-button">
                {questionIndex < questionsAmount ?
                (<button
                  className="shaped-button"
                  onClick={(e) => handleNextQuestion(e)}
                >
                  <img src={Triangle} alt="Go" className="shaped-button-icon" fill="var(--dark)"/>
                </button>)
                : (
                <button className="shaped-button"
                  onClick={() => finishSession(sessionStorage.getItem('sessionCode'))}
                >
                  <img src={Triangle} alt="Finish" className="shaped-button-icon" fill="var(--dark)"/>
                </button>)}
              </div>
            </div>
          </div>

          <div className="options-grid">
            {currentQuestion?.options.map((option, idx) => (
              <div key={idx} className={`option-item ${idx % 2 === 0 ? 'option-item-left' : 'option-item-right'}`}>
                <img src={questionOptions[idx]} alt={option.text} className={`option-image`} />
                <span className="option-text">{option.text}</span>
              </div>
            ))}
          </div>

          {/* <div className="process-footer">
            <div className="process-button-group">
              {questionIndex < questionsAmount && (
                <button onClick={handleNextQuestion} className="button">
                  Next
                </button>
              )}
              <span>Question: {questionIndex}/{questionsAmount}</span>
              <button onClick={() => finishSession(sessionStorage.getItem('sessionCode'))} className="nav-button">
                Finish
              </button>
            </div>
          </div> */}
        </div>
      )}
    </div>
  );
};

export default GameProcessAdmin;
