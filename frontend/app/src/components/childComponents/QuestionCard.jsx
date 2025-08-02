import React, { useState, useRef } from "react";
import styles from './QuestionCard.module.css';
import sampleImage from "../assets/sampleImage.png";
import { API_ENDPOINTS } from "../../constants/api";

const QuestionCard = ({ question, index, onChange }) => {
  const [showMenuIndex, setShowMenuIndex] = useState(null);
  const fileInputRef = useRef(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };
  
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
    setShowMenuIndex(null);
  };

  const handleDeleteOption = (i) => {
    const updatedOptions = question.options.filter((_, idx) => idx !== i);
    onChange({ ...question, options: updatedOptions });
    setShowMenuIndex(null);
  };

  const handleContextMenu = (e, i) => {
    e.preventDefault();
    setShowMenuIndex(i === showMenuIndex ? null : i);
  };

  const handleUploadImage = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await fetch(`${API_ENDPOINTS.QUIZ}/images/upload`, {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`Upload failed (${res.status}): ${errorText}`);
      }

      const data = await res.json();
      onChange({ ...question, image_url: data.url });
    } catch (err) {
      // console.error("Image upload failed:", err);
    }
  };

  return (
    <div className={styles['question-card']}>
      <div
        className={styles['image-upload-container']}
        onMouseEnter={() => setShowMenuIndex("image")}
        onClick={handleButtonClick}
      >
        <img
          src={question.image_url || sampleImage}
          alt="Question"
          className={styles['question-image']}
        />
        {showMenuIndex === "image" && (
          <div className={styles['upload-wrapper']}>
            <button className={styles['upload-button']}>Upload Image</button>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              onChange={handleUploadImage}
            />
          </div>
        )}
      </div>

      <div className={styles['question-content']}>
        <input
          type="text"
          placeholder="Enter question text"
          value={question.text}
          onChange={handleTextChange}
          className={styles['question-title-input']}
        />
        
        <div className={styles['answer-buttons']}>
          {question.options.map((answer, i) => (
            <div 
              key={i} 
              className={styles['answer-wrapper']} 
              onContextMenu={(e) => handleContextMenu(e, i)}
            >
              <button 
                className={`${styles['answer-button']} ${
                  answer.is_correct ? styles['highlight'] : ''
                }`}
              >
                <input
                  type="text"
                  className={styles['answer-input']}
                  placeholder="Enter answer"
                  value={answer.text}
                  onChange={(e) => handleAnswerChange(i, e.target.value)}
                />
              </button>

              {showMenuIndex === i && (
                <div className={styles['context-menu']}>
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