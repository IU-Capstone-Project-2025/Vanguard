// components/childComponents/QuestionCard.jsx
import React, { useState, useRef } from "react";
import "./QuestionCard.css";
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
  };

  const handleMarkIncorrect = (i) => {
    const updatedOptions = question.options.map((opt, idx) => ({
      ...opt,
      is_correct: idx !== i,
    }));
    onChange({ ...question, options: updatedOptions });
  }

  const handleDeleteOption = (i) => {
    const updatedOptions = question.options.filter((_, idx) => idx !== i);
    onChange({ ...question, options: updatedOptions });
  };

  const handleContextMenu = (e, i) => {
    e.preventDefault();
    setShowMenuIndex(i === showMenuIndex ? null : i);
  };

  const handleUploadImage = async (e) => {
    console.log("Uploading image...");
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("file", file); // название поля должно быть "image"

    try {
      const res = await fetch(`${API_ENDPOINTS.QUIZ}/images/upload/`, {
        method: "POST",
        body: formData,
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(`Upload failed (${res.status}): ${errorText}`);
      }

      const data = await res.json();
      const imageUrl = data.url;
      console.log(imageUrl)

      onChange({ ...question, image_url: imageUrl });
    } catch (err) {
      alert("Image upload failed: " + err.message);
    }
  };

  return (
    <div className="question-card">
      <div
        className="image-upload-container"
        onMouseEnter={() => setShowMenuIndex("image")}
        // onMouseLeave={() => setShowMenuIndex(null)}
        onClick={handleButtonClick}
      >
        <img
          src={question.image_url || sampleImage}
          alt="Question"
          className="question-image"
        />
        {showMenuIndex === "image" && (
          <div className="upload-wrapper">
            <button className="upload-button">Upload Image</button>
            <input
              ref={fileInputRef}
              id={`fileInput-${index}`}
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              onChange={(e) => {
                console.log("Image selected:", e.target.files[0]);
                handleUploadImage(e);
              }}
            />
          </div>
        )}

      </div>
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
