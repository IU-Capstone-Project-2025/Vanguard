import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './styles/FinalAdminPage.module.css';
import { API_ENDPOINTS } from '../constants/api';

const FinalAdminPage = () => {
  const [sessionCode] = useState(sessionStorage.getItem('sessionCode') || '');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const [players, setPlayers] = React.useState(new Map());
  
  const endSession = async (code) => {
    try {
      setIsLoading(true);
      const response = await fetch(`${API_ENDPOINTS.SESSION}/session/${code}/end`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      
      if (!response.ok) {
        throw new Error(`Failed to end session with code: ${code}`);
      }

      // Clear session storage
      sessionStorage.removeItem('sessionCode');
      sessionStorage.removeItem('quizData');
      sessionStorage.removeItem('currentQuestion');
      sessionStorage.removeItem('players');
      sessionStorage.removeItem('leaders');
      
      navigate('/');
    } catch (error) {
      console.error('Error ending the session:', error);
      setError('Failed to end session. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const storedLeaders = sessionStorage.getItem('leaders');
    const storedPlayers = sessionStorage.getItem('players');
    if (!storedLeaders || !storedPlayers) return;

    try {
      const leaders = JSON.parse(storedLeaders);
      const players = JSON.parse(storedPlayers);

      const nameMap = {};
      players.forEach(p => {
        if (p.name !== 'Admin') {
          nameMap[p.id] = p.name;
        }
      });

      const processed = leaders
        .filter(l => nameMap[l.user_id])
        .sort((a, b) => b.score - a.score)
        .slice(0, 3)
        .map((leader, index) => ({
          ...leader,
          name: nameMap[leader.user_id],
          position: index + 1,
        }));

      setPlayers(processed);
    } catch (err) {
      console.error('Error parsing leaders or players:', err);
      setError('Failed to load leaderboard data.');
    }
  }, []);

  return (
    <div className={styles['final-page-container']}>
      <h1 className={styles['final-title']}>Meet the Champions!</h1>
      
      {error && <div className={styles.error}>{error}</div>}
      
      <div className={styles.podium}>
        {players.length > 0 ? (
          players
            .sort((a, b) => a.position - b.position)
            .map((player) => (
              <div 
                key={player.user_id} 
                className={`${styles['podium-block']} ${styles[`position-${player.position}`]}`}
              >
                <span className={styles['player-name']}>{player.name}</span>
                <span className={styles['player-score']}>{player.score} pts</span>
              </div>
            ))
        ) : (
          <div className={styles['no-players']}>No players to display</div>
        )}
      </div>
      
      <div className={styles['buttons-container']}>
        <button 
          className={styles['final-primary-button']} 
          onClick={() => endSession(sessionCode)}
          disabled={isLoading}
        >
          {isLoading ? 'Ending Session...' : 'Back to Home'}
        </button>
      </div>
    </div>
  );
};

export default FinalAdminPage;