// components/childComponents/QuestionCard.jsx
import React, { useState } from "react";
import "./QuestionCard.css";
import sampleImage from "../assets/sampleImage.png";

const QuestionCard = ({ question, index, onChange }) => {
  const [showMenuIndex, setShowMenuIndex] = useState(null);

  const handleTextChange = (e) => {
    onChange({ ...question, text: e.target.value });
  };

  const handleAnswerChange = (i, newText) => {
    const updatedOptions = [...question.options];
    updatedOptions[i].text = newText;
    onChange({ ...question, options: updatedOptions });
  };

  const handleMarkCorrect = (i) => {
    const updatedOptions = question.options.map((opt, idx) => ({
      ...opt,
      is_correct: idx === i,
    }));
    onChange({ ...question, options: updatedOptions });
  };

  const handleDeleteOption = (i) => {
    const updatedOptions = question.options.filter((_, idx) => idx !== i);
    onChange({ ...question, options: updatedOptions });
  };

  const handleContextMenu = (e, i) => {
    e.preventDefault();
    setShowMenuIndex(i === showMenuIndex ? null : i);
  };

  return (
    <div className="question-card">
      <img src={sampleImage} alt="Question" className="question-image" />
      <div className="question-content">
        <input
          type="text"
          placeholder="Enter question text"
          value={question.text}
          onChange={handleTextChange}
          className="question-title-input"
        />
        <div className="answer-buttons">
          {question.options.map((answer, i) => (
            <div key={i} className="answer-wrapper" onContextMenu={(e) => handleContextMenu(e, i)}>
              <button className={`answer-button ${answer.is_correct ? "highlight" : ""}`}>
                <input
                  type="text"
                  className="answer-input"
                  placeholder="Enter answer"
                  value={answer.text}
                  onChange={(e) => handleAnswerChange(i, e.target.value)}
                />
              </button>

              {showMenuIndex === i && (
                <div className="context-menu">
                  <div onClick={() => handleMarkCorrect(i)}>Mark as correct</div>
                  <div onClick={() => handleDeleteOption(i)}>Delete</div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default QuestionCard;
