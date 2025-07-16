// ShowQuizStatistics.jsx
import React from 'react';
import { BarChart } from '@mui/x-charts/BarChart';
import { motion, AnimatePresence } from 'framer-motion';
import './ShowQuizStatistics.css'; // создайте по желанию для кастомизации

const ShowQuizStatistics = ({ stats, onClose }) => {
  const labels = Object.keys(stats || {});
  const values = Object.values(stats || {});

  return (
    <AnimatePresence>
      {stats && (
        <motion.div
          className="quiz-stats-overlay"
          initial={{ y: '-100vh' }}
          animate={{ y: 0 }}
          exit={{ y: '-100vh' }}
          transition={{ type: 'spring', stiffness: 80 }}
        >
          <h2 className="quiz-stats-title">Answer Statistics</h2>
          <div>
            <BarChart
              xAxis={[{ scaleType: 'band', data: labels }]}
              series={[{ data: values }]}
              width={500}
              height={300}
              colors={['#f9f3eb']}
            />
          </div>
          <button className="quiz-stats-button" onClick={onClose}>▶ Continue</button>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default ShowQuizStatistics;
