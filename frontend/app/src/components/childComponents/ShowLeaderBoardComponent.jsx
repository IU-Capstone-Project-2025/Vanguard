import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import styles from './ShowLeaderBoardComponent.module.css';

const ShowLeaderBoardComponent = ({ leaderboardData, onClose }) => {
  const [players, setPlayers] = React.useState(new Map());

  React.useEffect(() => {
    const stored = sessionStorage.getItem('players');
    const restored = stored ? new Map(JSON.parse(stored)) : new Map();
    setPlayers(restored);
  }, []);

  return (
    <AnimatePresence>
      {leaderboardData && (
        <motion.div
          className={styles["leaderboard-overlay"]}
          initial={{ y: '-100vh' }}
          animate={{ y: 0 }}
          exit={{ y: '-100vh' }}
          transition={{ type: 'spring', stiffness: 80 }}
        >
          <h1 className={styles["leaderboard-title"]}>Look! Here's our champions!</h1>
          {leaderboardData.users && (
            <div className={styles["leaderboard-list"]}>
              {leaderboardData.users.map((user, index) => (
                <div className={styles["leaderboard-row"]} key={user.user_id}>
                  <span className={styles["player-name"]}>
                    {players.get(user.user_id) || 'Unknown Player'}
                  </span>
                  <span className={styles["player-score"]}>{user.total_score}</span>
                </div>
              ))}
            </div>
          )}
          <button className={styles["leaderboard-next-button"]} onClick={onClose}>
            ▶ Next
          </button>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default ShowLeaderBoardComponent;
