import React, { useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import confetti from 'canvas-confetti';
import styles from './ShowQuizStatistics.module.css';
import PentagonYellow from '../assets/Pentagon-yellow.svg';
import CoronaIndigo from '../assets/Corona-indigo.svg';
import ArrowOrange from '../assets/Arrow-orange.svg';
import Cookie4Blue from '../assets/Cookie4-blue.svg';

const options = [PentagonYellow, CoronaIndigo, ArrowOrange, Cookie4Blue]

const ShowQuizStatistics = ({ stats, correct, onClose }) => {
  const headerRef = useRef(null);

  useEffect(() => {
    if (!headerRef.current || !correct) return;

    const rect = headerRef.current.getBoundingClientRect();
    const originX = (rect.left + rect.width / 2) / window.innerWidth;
    const originY = (rect.top + rect.height / 2) / window.innerHeight;

    confetti({
      particleCount: 100,
      spread: 80,
      origin: { x: originX, y: originY },
      colors: ['#00C851', '#33b5e5', '#ffbb33'],
    });
  }, [correct]);

  if (!stats) return null;

  const labels = Object.keys(stats);
  const values = Object.values(stats);
  const max = Math.max(...values);

  return (
    <AnimatePresence>
      <motion.div
        className={styles['quiz-stats-overlay']}
        initial={{ y: '-100vh' }}
        animate={{ y: 0 }}
        exit={{ y: '-100vh' }}
        transition={{ type: 'spring', stiffness: 80 }}
      >
        {/* <div
          ref={headerRef}
          className={`${styles['quiz-stats-header']} ${correct ? styles['correct'] : styles['incorrect']}`}
        >
          <h2 className={styles['quiz-stats-result']}>
            {correct ? 'Correct!' : 'Incorrect!'}
          </h2>
        </div>

        <h2 className={styles['quiz-stats-title']}>Answer Statistics</h2> */}

        <div className={styles['bars-container']}>
          {labels.map((label, idx) => {
            const heightPercent = Math.max((values[idx] / max) * 80, 10) + 20;

            return (
              <div className={styles['bar-wrapper']} key={label}>
                <motion.div
                  className={styles['bar']}
                  initial={{ height: 0 }}
                  animate={{ height: `${heightPercent}%` }}
                  transition={{ duration: 0.8, delay: idx * 0.4, ease: 'easeOut' }}
                >
                  <div className={styles['bar-icon']}>
                    <img src={options[idx]} alt='idx'/>
                  </div>
                </motion.div>
                <div className={styles['bar-feedback']}>
                  {[...Array(values[idx])].map((_, i) => (
                    <span key={i} className={i === 0 && correct ? styles['tick'] : styles['cross']}>
                      {i === 0 && correct ? '✔' : '✘'}
                    </span>
                  ))}
                </div>
              </div>
            );
          })}
        </div>

      </motion.div>
    </AnimatePresence>
  );
};

export default ShowQuizStatistics;
