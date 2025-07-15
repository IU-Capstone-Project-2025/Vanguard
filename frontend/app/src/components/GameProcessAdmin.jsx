import React, { useEffect, useState, useCallback } from 'react';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import './styles/GameProcess.css'
import { useNavigate } from 'react-router-dom';
import { API_ENDPOINTS } from '../constants/api';
import ShapedButton from './childComponents/ShapedButton';
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
  const questionsAmount = currentQuestion.questionsAmount - 1;
  console.log(`received questions amount [${questionsAmount}] and current index [${questionIndex}]`)
  const navigate = useNavigate();
  const [questionOptions, setQuestionOptions] = useState([
        Alien,
        Corona,
        Ghosty,
        Cookie6
      ] // Default empty question to avoid errors
    )

  useEffect(() => {
    const token = sessionStorage.getItem('jwt');
    if (!token) return;

    // console.log('Current question from sessionStorage:', currentQuestion);

    /* отписка при размонтировании */
    return () => {
      if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
      if (wsRefSession.current)  wsRefSession.current.onmessage  = null;
      if (wsRefRealtime.current) wsRefRealtime.current.onclose = {finishSession}
      if (wsRefSession.current) wsRefSession.current.onclose = {finishSession}
    };
  }, [connectSession, connectRealtime, wsRefSession, wsRefRealtime]);

const toNextQuestion = async (sessionCode) => {
    console.log("give the next question api call in game-process-admin");

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
      if (!response.ok) {
        throw new Error('Failed to start next question');
      }
      console.log('Next question started:', response);
    } catch (error) {
      console.error('Error starting next question:', error);
    }
  };


  const listenQuizQuestion = async (sessionCode) => {
    if (!sessionCode) {
      console.error('Session code is not available');
      return;
    }
    try {
      wsRefRealtime.current.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'question') {
          console.log('Received question:', data);
          let modifiedQuestion = {

          }
          
          setCurrentQuestion(data);
          sessionStorage.setItem('currentQuestion', JSON.stringify(data));
          return data
        }
      };
    } catch (error) {
      console.error('Error listening for quiz questions:', error);
      return
    }
  };

  const finishSession = async (code) => {
    try {
      const response = await fetch(`${API_ENDPOINTS.SESSION}/session/${code}/end`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error(`Failed to end session with code: ${code}`);
      }
      console.log(`end session with code: [${code}] response:`, response);
      // Очистка sessionStorage
      sessionStorage.removeItem('sessionCode');
      sessionStorage.removeItem('quizData');
      sessionStorage.removeItem('currentQuestion');
      // Закрытие WebSocket соединений
      closeWsRefRealtime();
      closeWsRefSession();
      navigate('/')
    } catch (error) {
      console.error('Error end the session:', error);
    }
    
  }

  /* -------- кнопка "Next" -------- */
  const handleNextQuestion = async (e) => {
    e.preventDefault();
    const sessionCode = sessionStorage.getItem('sessionCode');

    // if (questionIndex === questionsAmount) {
    //   await finishSession(sessionCode)
    // }

    toNextQuestion(sessionCode);
    setQuestionIndex((prevIndex) => prevIndex + 1);

    const quizData = await listenQuizQuestion(sessionCode)
  };

  /* -------- UI -------- */
  return (
    <div className="game-process">
      <div className='controller-question-title'>
        <h1>Live Quiz</h1>

        <p>Question {questionIndex + 1}</p>
      
        <h2>{currentQuestion ? currentQuestion.text : 'Waiting for question…'}</h2>
      </div>

      <div className="options-grid">
        {currentQuestion &&
          currentQuestion.options.map((option, idx) => (
            <ShapedButton 
              key={idx}
              shape={questionOptions[idx]}
              text={option.text} 
              onClick={
                () => {console.log('svg clicked')}
              }
            />
          ))
        }
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
    </div>
  );
};

export default GameProcessAdmin;
