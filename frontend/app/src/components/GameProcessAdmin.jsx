import React, { useEffect, useState } from 'react';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import styles from './styles/GameProcess.module.css';
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../constants/api';
import ShowLeaderBoardComponent from './childComponents/ShowLeaderBoardComponent';
import PentagonYellow from './assets/Pentagon-yellow.svg';
import CoronaIndigo from './assets/Corona-indigo.svg';
import ArrowOrange from './assets/Arrow-orange.svg';
import Cookie4Blue from './assets/Cookie4-blue.svg';
import Triangle from './assets/Triangle.svg';

const GameProcessAdmin = () => {
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const [currentQuestion, setCurrentQuestion] = useState(
    sessionStorage.getItem('currentQuestion') 
      ? JSON.parse(sessionStorage.getItem('currentQuestion')) 
      : {}
  );
  const [questionIndex, setQuestionIndex] = useState(1);
  const [leaderboardVisible, setLeaderboardVisible] = useState(false);
  const [leaderboardData, setLeaderboardData] = useState(null);
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [answerCounter, setAnswerCounter] = useState(0);
  const [playersAmount, setPlayersAmount] = useState(0);

  const questionsAmount = currentQuestion.questionsAmount || 0;
  const navigate = useNavigate();
  const questionOptions = [PentagonYellow, CoronaIndigo, ArrowOrange, Cookie4Blue];

  useEffect(() => {
    const token = sessionStorage.getItem('jwt');
    const sessionCode = sessionStorage.getItem('sessionCode');
    
    if (!token || !sessionCode) {
      setError('Missing session credentials');
      return;
    }

    const stored = sessionStorage.getItem('players');
    const restored = stored ? new Array(JSON.parse(stored)) : [];
    
    setPlayersAmount(restored.length)
    
    try {
      connectSession(token, sessionCode);
      connectRealtime(token, sessionCode);
    } catch (err) {
      setError('Failed to connect to game server');
      console.error('Connection error:', err);
    }

    const handleRealtimeMessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Realtime message:', data);

        if (data.type === 'leaderboard') {
          sessionStorage.setItem('leaders', JSON.stringify(data.payload.users));
          setLeaderboardData(data.payload);
          setLeaderboardVisible(true);
        } else if (data.type === 'question') {
          setAnswerCounter(0);
          setCurrentQuestion(data);
          setQuestionIndex(data.questionId);
          sessionStorage.setItem('currentQuestion', JSON.stringify(data));
        } else if (data.type === 'user_answered') {
          setAnswerCounter((prev) => prev + 1)
        }

      } catch (err) {
        console.error('Error processing message:', err);
        setError('Error processing game data');
      }
    };

    if (wsRefRealtime.current) {
      wsRefRealtime.current.onmessage = handleRealtimeMessage;
    }

    return () => {
      if (wsRefRealtime.current) {
        wsRefRealtime.current.onmessage = null;
      }
      if (wsRefSession.current) {
        wsRefSession.current.onmessage = null;
      }
    };
  }, [connectSession, connectRealtime, wsRefSession, wsRefRealtime, playersAmount]);

  const toNextQuestion = async (sessionCode) => {
    if (!sessionCode) {
      setError('Session code is missing');
      return;
    }

    try {
      setIsLoading(true);
      const response = await fetch(`${API_ENDPOINTS.SESSION}/session/${sessionCode}/nextQuestion`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      
      if (!response.ok) {
        throw new Error('Failed to start next question');
      }
    } catch (error) {
      console.error('Error starting next question:', error);
      setError('Failed to advance to next question');
    } finally {
      setIsLoading(false);
    }
  };

  const finishSession = async () => {
    try {
      setIsLoading(true);

      const sessionCode = sessionStorage.getItem('sessionCode');
      await toNextQuestion(sessionCode);

      sessionStorage.removeItem('quizData');
      sessionStorage.removeItem('currentQuestion');
      
      closeWsRefRealtime();
      closeWsRefSession();
      
      navigate('/final');
    } catch (error) {
      console.error('Error finishing session:', error);
      setError('Failed to end session properly');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLeaderboardClick = () => {
    if (!wsRefRealtime.current) {
      setError('Connection to game server lost');
      return;
    }

    try {
      wsRefRealtime.current.send(JSON.stringify({ type: 'next_question' }));
      setLeaderboardVisible(false);
    } catch (err) {
      console.error('Error sending next question:', err);
      setError('Failed to proceed to next question');
    }
  };

  const handleNextQuestion = async (e) => {
    e.preventDefault();
    await toNextQuestion(sessionStorage.getItem('sessionCode'));
  };

  return (
    <div className={styles['game-process']}>
      {error && <div className={styles.error}>{error}</div>}

      {leaderboardVisible && leaderboardData ? (
        <ShowLeaderBoardComponent
          leaderboardData={leaderboardData}
          onClose={handleLeaderboardClick} 
        />
      ) : (
        <div className={styles['game-container']}>
          <div className={styles['question-section']}>
            <div className={styles['question-header']}>
              <div className={styles['answer-indicator']}>
                <span className={styles['indicator-text']}>{answerCounter}/{playersAmount}</span>
              </div>

              {currentQuestion.payload && (
                <img
                  src={currentQuestion.payload}
                  alt="Question"
                  className={styles['process-question-image']}
                />
              )}

              <div className={styles['timer-indicator']}>
                <span className={styles['indicator-text']}>sample</span>
              </div>
            </div>

            <div className={styles['question-body']}>
              <div className={styles['question-number-bubble']}>
                {questionIndex}
              </div>
              <h2 className={styles['question-text']}>
                {currentQuestion?.text || 'Waiting for question...'}
              </h2>
              
              <div className={styles['question-go-button']}>
                <button
                  className={styles['shaped-button']}
                  onClick={questionIndex < questionsAmount ? handleNextQuestion : finishSession}
                  disabled={isLoading}
                >
                  <img 
                    src={Triangle} 
                    alt={questionIndex < questionsAmount ? "Next" : "Finish"} 
                    className={styles['shaped-button-icon']} 
                  />
                </button>
              </div>
            </div>
          </div>

          <div className={styles['options-grid']}>
            {currentQuestion?.options?.map((option, idx) => (
              <div 
                key={idx} 
                className={`${styles['option-item']} ${
                  idx % 2 === 0 ? styles['option-item-left'] : styles['option-item-right']
                }`}
              >
                <img 
                  src={questionOptions[idx]} 
                  alt={option.text} 
                  className={styles['option-image']} 
                />
                <span className={styles['option-text']}>{option.text}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default GameProcessAdmin;