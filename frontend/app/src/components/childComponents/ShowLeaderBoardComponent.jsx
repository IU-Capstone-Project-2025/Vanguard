import React, { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import styles from './ShowLeaderBoardComponent.module.css';
import ArrowUp from '../assets/Arrow-up.svg';
import ArrowDown from '../assets/Arrow-down.svg';

const ShowLeaderBoardComponent = ({ leaderboardData, onClose }) => {
  const [players, setPlayers] = useState(new Map());

  useEffect(() => {
    const stored = sessionStorage.getItem('players');
    const restored = stored ? new Map(JSON.parse(stored)) : new Map();
    setPlayers(restored);
  }, []);

  return (
    <AnimatePresence>
      {leaderboardData && (
        <motion.div
          className={styles.overlay}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
        >
          <div className={styles.container}>
            <motion.h1
              className={styles.title}
              initial={{ y: '-200px', opacity: 0 }}
              animate={{ y: '0px', opacity: 1 }}
              transition={{ delay: 0.2, type: 'spring', stiffness: 80 }}
            >
              Scoreboard!
            </motion.h1>

            <motion.button
              className={styles.next}
              onClick={onClose}
              initial={{ y: '-320px', x: '560px', opacity: 0 }}
              animate={{ y: '-120px', x: '560px', opacity: 1 }}
              transition={{ delay: 0.2, type: 'spring', stiffness: 80 }}
            >
              ▶ Next
            </motion.button>

            <motion.div
              className={styles.board}
              initial={{ y: '100vh' }}
              animate={{ y: 0 }}
              exit={{ y: '100vh' }}
              transition={{ type: 'spring', stiffness: 80 }}
            >
              {leaderboardData.users
                ?.filter((user) => players.has(user.user_id))
                .map((user) => (
                  <div className={styles.row} key={user.user_id}>
                    <span className={styles.place}>{user.place}</span>
                    <div className={styles['user-info']}>
                      <span className={styles.name}>
                        {players.get(user.user_id)}
                      </span>
                      <span className={styles.score}>{user.total_score}</span>
                      {user.progress === "Up" ? (
                        <img src={ArrowUp} alt="Up" className={styles.arrow} />
                      ) : user.progress === "Down" ? (
                        <img src={ArrowDown} alt="Down" className={styles.arrow} />
                      ) : (
                        <div className={styles.arrow}></div>
                      )}
                    </div>
                  </div>
              ))}
            </motion.div>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default ShowLeaderBoardComponent;
