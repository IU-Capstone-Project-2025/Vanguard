import React, { useEffect, useState, useCallback } from 'react';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import './styles/GameProcess.css'
import { useNavigate } from 'react-router-dom';

  const GameProcessAdmin = () => {
    const { wsRefSession, connectSession } = useSessionSocket();
    const { wsRefRealtime, connectRealtime } = useRealtimeSocket();
    const [currentQuestion, setCurrentQuestion] = useState(sessionStorage.getItem('currentQuestion') != undefined ?
    JSON.parse(sessionStorage.getItem('currentQuestion')) : {});
    const [questionIndex, setQuestionIndex] = useState(0);
    const questionsAmount = useState(currentQuestion.quiestionsAmount - 1);
    const navigate = useNavigate()
    
    useEffect(() => {
      const token = sessionStorage.getItem('jwt');
      if (!token) return;

      // console.log('Current question from sessionStorage:', currentQuestion);

      /* отписка при размонтировании */
      return () => {
        if (wsRefRealtime.current) wsRefRealtime.current.onmessage = null;
        if (wsRefSession.current)  wsRefSession.current.onmessage  = null;
      };
    }, [connectSession, connectRealtime, wsRefSession, wsRefRealtime]);

  const toNextQuestion = async (sessionCode) => {
      console.log("give the next question api call in game-process-admin");

      if (!sessionCode) {
        console.error('Session code is not available');
        return;
      }
      try {
        const response = await fetch(`/api/session/session/${sessionCode}/nextQuestion`, {
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
        const response = await fetch(`/api/session/session/${code}/end`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        if (!response.ok) {
          throw new Error(`Failed to end session with code: ${code}`);
        }
        console.log(`end session with code: [${code}] response:`, response);
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
        <h1>Live Quiz</h1>

        <p>Question {questionIndex + 1}</p>

        <div className="question-block">
          <h2>{currentQuestion ? currentQuestion.text : 'Waiting for question…'}</h2>
        </div>

        <div className="options-grid">
          {currentQuestion &&
            currentQuestion.options.map((option, idx) => (
              <button key={idx} className="option-button">
                {option.text}
              </button>
            ))}
        </div>

        <div className="navigation-buttons">
          <button onClick={handleNextQuestion} className="nav-button">
            {/* { (questionIndex < questionsAmount) ? "Next" : "Finish" } */}
            Next
          </button>
        </div>
      </div>
    );
  };

  export default GameProcessAdmin;
 