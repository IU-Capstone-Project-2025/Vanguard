import React from "react";
import styles from './QuizPreviewModal.module.css';
import { motion, AnimatePresence } from "framer-motion";
import sampleImage from "../assets/sampleImage.png";

const backdropVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

const modalVariants = {
  hidden: { y: `${-100}px`, x: `${0}px`, scale: 0, opacity: 0 },
  visible: {
    y: "0",
    x: "0",
    scale: 1,
    opacity: 1,
    transition: { type: "spring", damping: 20, stiffness: 120 },
  },
  exit: {
    y: "100vh",
    opacity: 0,
    transition: { duration: 0.3 },
  },
};

const QuizPreviewModal = ({ quiz, onClose, coordinates, onStart, onEdit }) => {
  if (!quiz) return null;

  return (
    <AnimatePresence>
      <motion.div
        className={styles['modal-overlay']}
        variants={backdropVariants}
        initial="hidden"
        animate="visible"
        exit="hidden"
        onClick={onClose}
      >
        <motion.div
          className={styles['modal-wrapper']}
          variants={modalVariants}
          initial="hidden"
          animate="visible"
          exit="exit"
          onClick={(e) => e.stopPropagation()}
        >
          <div className={styles['modal-navbar']}>
            <button onClick={onClose}>⤫</button>
          </div>

          <div className={styles['modal-title']}>
            <h2>{quiz.title}</h2>
            <p>{quiz.questions.length} questions</p>
          </div>

          <div className={styles['questions-list']}>
            {quiz.questions.map((q, index) => (
              <div key={index} className={styles['question-preview']}>
                <img
                  src={q.image_url || sampleImage}
                  alt="Question"
                  className={styles['question-image']}
                />
                <div className={styles['question-details']}>
                  <p>{q.text}</p>
                  <div className={styles['answers']}>
                    {q.options.map((ans, i) => (
                      <button
                        key={i}
                        className={`${styles['answer-button']} ${
                          ans.is_correct ? styles['correct'] : ''
                        }`}
                      >
                        {ans.text}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>

          <div className={styles['modal-footer']}>
            <button 
              className={styles['secondary-button']} 
              onClick={onEdit}
            >
              Edit quiz
            </button>
            <button 
              className={styles['primary-button']} 
              onClick={onStart}
            >
              Start Quiz
            </button>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
};

export default QuizPreviewModal;