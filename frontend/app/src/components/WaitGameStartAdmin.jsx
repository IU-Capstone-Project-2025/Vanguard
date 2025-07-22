import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSessionSocket } from '../contexts/SessionWebSocketContext';
import { useRealtimeSocket } from '../contexts/RealtimeWebSocketContext';
import styles from './styles/WaitGameStartAdmin.module.css';
import { API_ENDPOINTS, BASE_URL } from '../constants/api';
import QRCodeModal from './childComponents/QRCodeModal.jsx'; // путь подкорректируй


const WaitGameStartAdmin = () => {
  const navigate = useNavigate();
  const { wsRefSession, connectSession, closeWsRefSession } = useSessionSocket();
  const { connectRealtime, wsRefRealtime, closeWsRefRealtime } = useRealtimeSocket();
  const [sessionCode] = useState(sessionStorage.getItem('sessionCode'));
  const [players, setPlayers] = useState(new Map());
  const [isStarting, setIsStarting] = useState(false);
  const [showQRModal, setShowQRModal] = useState(false);


  const handleSessionMessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      const newPlayers = new Map();
      for (const [userId, name] of Object.entries(data)) {
        newPlayers.set(userId, name);
      }
      sessionStorage.setItem('players', JSON.stringify([...newPlayers]));
      setPlayers(newPlayers);
    } catch (e) {
      console.error('Invalid session message:', e);
    }
  };

  const handleRealtimeMessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.type === 'question') {
        sessionStorage.setItem('currentQuestion', JSON.stringify(data));
      }
    } catch (e) {
      console.error('Failed to parse realtime message:', e);
    }
  };

  useEffect(() => {
    const token = sessionStorage.getItem('jwt');
    if (!token) return;

    connectSession(token, handleSessionMessage);
    connectRealtime(token, handleRealtimeMessage);

    return () => {
      closeWsRefRealtime();
      closeWsRefSession();
    };
  }, [connectSession, connectRealtime]);

  const finishSession = async () => {
    try {
      await fetch(`${API_ENDPOINTS.SESSION}/session/${sessionCode}/end`, {
        method: 'POST',
      });
      sessionStorage.removeItem('sessionCode');
      sessionStorage.removeItem('quizData');
      sessionStorage.removeItem('currentQuestion');
      navigate('/');
    } catch (error) {
      console.error('Error ending session:', error);
    }
  };

  const handleKick = async (userId) => {
    try {
      await fetch(`${API_ENDPOINTS.SESSION}/delete-user?code=${sessionCode}&userId=${userId}`, {
        method: 'GET',
      });
    } catch (e) {
      console.error("Error kicking player:", e);
    }
  };

  const startGame = async () => {
    setIsStarting(true);
    try {
      await fetch(`${API_ENDPOINTS.SESSION}/session/${sessionCode}/nextQuestion`, {
        method: 'POST',
      });
      navigate(`/game-controller/${sessionCode}`);
    } catch (error) {
      console.error('Error starting game:', error);
      setIsStarting(false);
    }
  };

  return (
    <div className={styles['wait-admin-wrapper']}>
      {showQRModal && (
        <QRCodeModal
          code={`${BASE_URL.REACT_APP_BASE_URL}/join/${sessionCode}`} // здесь вставь нужный URL
          onClose={() => setShowQRModal(false)}
        />
      )}
      <h1 className={styles['waiting-title']}>
        Waiting for your awesome <span className={styles['highlight']}>crew</span>...
      </h1>
      
      <div className={styles['wait-layout']}>
        <div className={styles['players-grid']}>
          {[...players.entries()].map(([id, name]) => (
            <div key={id} className={styles['player-box']} onClick={() => handleKick(id)}>
              <span style={{ '--name-length': name.length }}>{name}</span>
              <div className={styles['tooltip']}>Click to kick user</div>
            </div>
          ))}
        </div>

        <div className={styles['right-panel']}>
          <div className={styles['session-box']}>
            <div className={styles['session-code']}># <strong>{sessionCode}</strong></div>
            <div className={styles['session-count']}>👤 {players.size}/40</div>
            <button onClick={() => setShowQRModal(true)} className={styles['start-btn']}>
              📱 Show QR
            </button>
            <button onClick={finishSession} className={styles['terminate-btn']}>
              ✖ Terminate
            </button>
            <button 
              onClick={startGame} 
              disabled={isStarting || players.size === 0}
              className={styles['start-btn']}
            >
              {isStarting ? 'Starting...' : '▶ Start'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default WaitGameStartAdmin;