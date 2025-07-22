import React, { useEffect, useRef } from 'react';
import { BarChart } from '@mui/x-charts/BarChart';
import { motion, AnimatePresence } from 'framer-motion';
import confetti from 'canvas-confetti';
import styles from './ShowQuizStatistics.module.css';

const ShowQuizStatistics = ({ stats, correct, onClose }) => {
  const chartRef = useRef(null);
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

  return (
    <AnimatePresence>
      <motion.div
        className={styles['quiz-stats-overlay']}
        initial={{ y: '-100vh' }}
        animate={{ y: 0 }}
        exit={{ y: '-100vh' }}
        transition={{ type: 'spring', stiffness: 80 }}
      >
        <div
          ref={headerRef}
          className={`${styles['quiz-stats-header']} ${
            correct ? styles['correct'] : styles['incorrect']
          }`}
        >
          <h2 className={styles['quiz-stats-result']}>
            {correct ? 'Correct!' : 'Incorrect!'}
          </h2>
        </div>

        <h2 className={styles['quiz-stats-title']}>Answer Statistics</h2>
        
        <div ref={chartRef} className={styles['chart-container']}>
          <BarChart
            xAxis={[{ scaleType: 'band', data: labels }]}
            series={[{ data: values }]}
            width={500}
            height={300}
            colors={['#f9f3eb']}
          />
        </div>
      </motion.div>
    </AnimatePresence>
  );
};

export default ShowQuizStatistics;