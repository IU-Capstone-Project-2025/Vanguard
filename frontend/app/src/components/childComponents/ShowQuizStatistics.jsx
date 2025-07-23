import React, { useEffect, useRef } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import confetti from 'canvas-confetti';
import styles from './ShowQuizStatistics.module.css';
import PentagonYellow from '../assets/Pentagon-yellow.svg';
import CoronaIndigo from '../assets/Corona-indigo.svg';
import ArrowOrange from '../assets/Arrow-orange.svg';
import Cookie4Blue from '../assets/Cookie4-blue.svg';

const icons = [PentagonYellow, CoronaIndigo, ArrowOrange, Cookie4Blue];

const ShowQuizStatistics = ({ stats, correct, onClose, options }) => {
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

  if (!stats || !options) return null;

  const keys = Object.keys(stats);
  const values = Object.values(stats);
  const max = Math.max(...values, 1); // защита от деления на 0

  return (
    <AnimatePresence>
      <motion.div
        className={styles['quiz-stats-overlay']}
        initial={{ y: '-100vh' }}
        animate={{ y: 0 }}
        exit={{ y: '-100vh' }}
        transition={{ type: 'spring', stiffness: 80 }}
      >
        <div className={styles['bars-container']}>
          {keys.map((key, idx) => {
            const count = stats[key];
            const isCorrect = options?.[idx]?.is_correct;
            const heightPercent = Math.max((count / max) * 70, 25); // минимум 25% для визуала

            return (
              <div className={styles['bar-wrapper']} key={key}>
                <motion.div
                  className={styles['bar']}
                  initial={{ height: 0 }}
                  animate={{ height: `${heightPercent}%` }}
                  transition={{ duration: 0.8, delay: idx * 0.3 }}
                >
                  <div className={styles['bar-icon']}>
                    <img src={icons[idx % icons.length]} alt="option-icon" />
                  </div>

                  <div className={styles['bar-info']}>
                    {isCorrect && <span className={styles['tick']}>✔</span>}
                    <span className={styles['bar-count']}>{count}</span>
                  </div>
                </motion.div>
              </div>
            );
          })}
        </div>
      </motion.div>
    </AnimatePresence>
  );
};

export default ShowQuizStatistics;
