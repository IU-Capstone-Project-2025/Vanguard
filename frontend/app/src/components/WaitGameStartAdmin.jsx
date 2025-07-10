import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import './styles/WaitGameStartAdmin.css';
import { API_ENDPOINTS } from '../constants/api';

const WaitGameStartAdmin = () => {
  const navigate = useNavigate();
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { connectRealtime, wsRefRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const [sessionCode, setSessionCode] = useState(sessionStorage.getItem('sessionCode') || null);
  const [players, setPlayers] = useState(new Map());
  const [hasClickedNext, setHasClickedNext] = useState(false)

  const extractPlayersFromMessage = (data) => {
    setPlayers((prevPlayers) => {
      const newPlayers = new Map(prevPlayers)
      for (const [userId,name] of Object.entries(data)) {
        if (!newPlayers.has(userId)) {
          newPlayers.set(userId, name);
        }
      }
      return newPlayers;
    });
  };

  useEffect(() => {
    const token = sessionStorage.getItem('jwt');
    if (!token) return;

    connectSession(token, (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log('Received session message:', data);
        extractPlayersFromMessage(data);
      } catch (e) {
        console.error('⚠️ Invalid session WS message:', event.data);
      }
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

    wsRefRealtime.current.onclose = () => {
      closeConnection();
    }
    wsRefSession.current.onclose = () => {
      closeConnection();
    }
  }, [connectSession, connectRealtime]);

  const closeConnection = () => {
    closeWsRefRealtime();
    closeWsRefSession();
    navigate('/');
  }

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
      closeConnection();
    } catch (error) {
      console.error('Error end the session:', error);
    }
    
  }

  const handleKick = async (idToRemove) => {
    console.log(`Kick user with id [${idToRemove}]`)
    try {
      const queryParams = new URLSearchParams(
        {
          code: sessionCode,
          userId: idToRemove
        }
      );
      const response = await fetch(`${API_ENDPOINTS.SESSION}/delete-user?${queryParams}`,
        {
          method: 'GET',
          'Content-Type': 'application/json'
        }
      )
      if (response.status !== 200) {
        throw new Error(`Failed to kick player with id: ${idToRemove}`);
      }
      console.log('Kicked player. response: ', response)
    }  
    catch (e) {
      console.error("Error with kicking: ", e)
    }
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
    e.preventDefault();
    setHasClickedNext(true);
    const sessionCode = sessionStorage.getItem('sessionCode');
    await toNextQuestion(sessionCode);
    navigate(`/game-controller/${sessionCode}`);
  };

  const handleTerminate = () => {
    finishSession(sessionCode);
  };

  return (
    <div className="wait-admin-container">
      <div className="wait-admin-panel">
        <h1>Now let's wait your friends <br /> Code: #{sessionCode}</h1>
        <div className="admin-button-group">
          <button onClick={handleStart} disabled={hasClickedNext}>▶️ Start</button>
          <button onClick={handleTerminate}>❌ Terminate</button>
        </div>
        <div className="players-grid">
          {(Array.from(players.entries())).map(([id, name]) => (
            <div key={id + name} className="player-box">
              <span>{name}</span>
              {name != "Admin" &&
              <button
                className="kick-button"
                onClick={() => handleKick(id)}
                title={`Kick ${name}`}
              >
                ❌
              </button>}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WaitGameStartAdmin;
