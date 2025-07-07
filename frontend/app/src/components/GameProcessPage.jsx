import React from "react";
import "./styles/GameProcess.css";

const GameProcess = ({ question, options, onAnswer }) => {
  return (
    <div className="game-process-container">
      <div className="game-process-panel">
        <p>And choose the correct answer</p>
        <div className="options-grid">
          {options.map((option, index) => (
            <button
              key={index}
              className="option-button"
              onClick={() => onAnswer(option)}
            >
              {option}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default GameProcess;
