import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './styles/FinalAdminPage.module.css';
import { API_ENDPOINTS } from '../constants/api';

const FinalAdminPage = () => {
  const [sessionCode] = useState(sessionStorage.getItem('sessionCode') || '');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const [visibleBlocks, setVisibleBlocks] = useState([]);

  const [leaders, setLeaders] = useState(new Map());
  const [players, setPlayers] = useState(new Map())

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
    const storedPlayersRaw = sessionStorage.getItem('players');
    if (!storedLeaders || !storedPlayersRaw) return;

    try {
      const leaders = JSON.parse(storedLeaders);
      const playersEntries = JSON.parse(storedPlayersRaw);
      const playersMap = new Map(playersEntries);

      setPlayers(playersMap)
      console.log(playersMap)

      const processed = leaders
        .filter(l => playersMap.get(l.user_id) && playersMap.get(l.user_id) !== 'Admin')
        .sort((a, b) => b.total_score - a.total_score)
        // .slice(0, 3)
        .map((leader, index) => [
          leader.user_id,
          {
            name: playersMap.get(leader.user_id),
            score: leader.total_score,
            position: index + 1,
          }
        ]);

      setLeaders(new Map(processed));
      console.log(leaders)

      // Анимации по таймеру
      processed.forEach(([, leader], i) => {
        setTimeout(() => {
          setVisibleBlocks(prev => [...prev, leader.position]);
        }, i * 700); // задержка 700мс между блоками
      });

    } catch (err) {
      console.error('Error parsing leaders or players:', err);
      setError('Failed to load leaderboard data.');
    }
  }, [setLeaders]);

  return (
    <div className={styles['final-page-container']}>
      <h1 className={styles['final-title']}>
        Say Hello to the <span className={styles.highlight}>Heroes</span>!
      </h1>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.podium}>
        {[...leaders.values()].length > 0 ? (
          [2, 1, 3] // порядок отображения: сначала 2-е, потом 1-е, потом 3-е
            .map((place) => {
              const player = [...leaders.values()].find(p => p.position === place);
              return (
                player &&
                visibleBlocks.includes(player.position) && (
                  <div
                    key={player.position}
                    className={`${styles['podium-block']} ${styles[`position-${player.position}`]} ${styles['animated']}`}
                  >
                    <span className={styles['player-name']}>{player.name}</span>
                    <span className={styles['player-place']}>{player.position}</span>
                  </div>
                )
              );
            })
        ) : (
          <div className={styles['no-players']}>No players to display</div>
        )}
      </div>

      <div className={styles.leaderboard}>
        {[...leaders.values()].map((player, idx) => (
          <div className={styles['leaderboard-row']} key={player.name + idx}>
            <span className={styles.place}>{player.position}</span>
            <div className={styles['user-info']}>
              <span className={styles.name}>{player.name}</span>
              <span className={styles.score}>{player.score}</span>
            </div>
          </div>
        ))}
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
