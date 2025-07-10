import React from "react";
import "./QuizPreviewModal.css";
import { motion, AnimatePresence, scale } from "framer-motion";
import sampleImage from "../assets/sampleImage.png";

const backdropVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

const modalVariants = {
//   hidden: { y: "-100vh", opacity: 0 },
    hidden: {y: "-100vh", scale: 0, opacity: 0},
  visible: {
    y: "0",
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
    const x = coordinates[0]
    const y = coordinates[1]

    const modalVariants = {
    //   hidden: { y: "-100vh", opacity: 0 },
        hidden: {y: `${-100}px`, x: `${0}px`, scale: 0, opacity: 0},
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

  if (!quiz) return null;

  const handleEdit = (e) => {
    e.preventDefault();
    onEdit(quiz);
  }

  const handleStart = () => {
    onStart(quiz);
  }

  return (
    <AnimatePresence>
      <motion.div
        className="modal-overlay"
        variants={backdropVariants}
        initial="hidden"
        animate="visible"
        exit="hidden"
        onClick={onClose}
      >
        <motion.div
          className="modal-wrapper"
          variants={modalVariants}
          initial="hidden"
          animate="visible"
          exit="exit"
          onClick={(e) => e.stopPropagation()}
        >
          <div className="modal-navbar">
            <button onClick={onClose}>⤫</button>
          </div>

          <div className="modal-title">
            <h2>{quiz.title}</h2>
            <p>{quiz.questions.length} questions</p>
          </div>

          {quiz.questions.map((q, index) => (
            <div key={index} className="question-preview">
                <img
                    src={q.imageURL || sampleImage}
                    alt="Question"
                />
                <div className="question-details">
                    <p>{q.text}</p>
                    <div className="answers">
                    {q.options.map((ans, i) => (
                        <button
                        key={i}
                        className={`answer-button ${
                            ans.is_correct ? "correct" : ""
                        }`}
                        >
                        {ans.text}
                        </button>
                    ))}
                    </div>
                </div>
            </div>
          ))}

          <div className="modal-footer">
            <button className="secondary-button" disabled >Edit quiz</button>
            <button className="primary-button" disabled >Start Quiz</button>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>
  );
};

export default QuizPreviewModal;
