import React, { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import './styles/WaitGameStartAdmin.css';
import { API_ENDPOINTS } from '../constants/api';


const WaitGameStartAdmin = () => {
  const navigate = useNavigate();
  const { wsRefSession, connectSession } = useSessionSocket();
  const { wsRefRealtime, connectRealtime } = useRealtimeSocket();
  const [sessionCode, setSessionCode] = useState(sessionStorage.getItem('sessionCode') || null);
  const [players, setPlayers] = useState([]);
  const [hasClickedNext, setHasClickedNext] = useState(false)

  const extractPlayersFromMessage = (data) => {
    if (!Array.isArray(data)) return;
    setPlayers((prevPlayers) => {
      const newPlayers = data.map((name, index) => ({
        id: prevPlayers.length + index + 1,
        name,
      }));
      return [...prevPlayers, ...newPlayers];
    });
  };

  useEffect(() => {
    const token = sessionStorage.getItem('jwt');
    if (!token) return;

    connectSession(token, (event) => {
      try {
        const data = JSON.parse(event.data);
        if (Array.isArray(data)) {
          setPlayers(data.map((name, i) => ({ id: i + 1, name })));
        }
      } catch (e) {
        console.error('⚠️ Invalid session WS message:', event.data);
      }
      console.log('Received session message:', event.data);
    });

    connectRealtime(token, (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'question') {
          console.log('Got question:', data);
          sessionStorage.setItem('currentQuestion', JSON.stringify(data));
        }
      } catch (e) {
        console.error('⚠️ Failed to parse realtime WS message:', event.data);
      }
    });
  }, [connectSession, connectRealtime]);

  const handleKick = (idToRemove) => {
    setPlayers((prev) => prev.filter((player) => player.id !== idToRemove));
    // TODO: отправить на backend сигнал о кике игрока по id
  };

  const toNextQuestion = async (sessionCode) => {
    console.log("give the next question api call in wait-admin");
    if (!sessionCode) {
      console.error('Session code is not available');
      return;
    }
    try {
      const response = await fetch(`${API_ENDPOINTS.SESSION}/session/${sessionCode}/nextQuestion`, {
        method: 'POST',
      });
      if (response.status !== 200) {
        throw new Error('Failed to start next question');
      }
      console.log('Next question started', response);
    } catch (error) {
      console.error('Error starting next question:', error);
    }
  };

  const handleStart = async (e) => {
    e.preventDefault()
    setHasClickedNext(true)
    const sessionCode = sessionStorage.getItem('sessionCode');
    await toNextQuestion(sessionCode);

    // const quizData = await listenQuizQuestion(sessionCode, wsRefRealtime);
    // sessionStorage.setItem('currentQuestion', JSON.stringify(quizData));

    navigate(`/game-controller/${sessionCode}`);
  };

  const handleTerminate = () => {
    navigate('/');
  };

  return (
    <div className="wait-admin-container">
      <div className="wait-admin-panel">
        <h1>Now let's wait your friends <br/> Code: #{sessionCode}</h1>
        <div className="admin-button-group">
          <button onClick={(e)=>{handleStart(e);}} disabled={hasClickedNext}>▶️ Start</button>
          <button onClick={handleTerminate}>▶️ Terminate</button>
        </div>
        <div className="players-grid">
          {players.map((player) => (
            <div key={player.id} className="player-box">
              <span>{player.name}</span>
              <button
                className="kick-button"
                onClick={() => handleKick(player.id)}
                title={`Kick ${player.name}`}
              >
                ❌
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WaitGameStartAdmin;
