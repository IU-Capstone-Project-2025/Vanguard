import React, { useState } from "react";
import QuestionCard from "./childComponents/QuestionCard";
import styles from './styles/ConstructorPage.module.css';
import { useNavigate } from "react-router-dom";
import { API_ENDPOINTS } from "../constants/api.js";

const ConstructorPage = () => {
  const navigate = useNavigate();
  const [quiz, setQuiz] = useState({
    title: "",
    description: "",
    is_public: true,
    tags: [],
    questions: [],
  });

  const handleAddQuestion = () => {
    const newQuestion = {
      type: "single_choice",
      text: "",
      image_url: "",
      time_limit: 5,
      options: [
        { text: "", image_url: "", is_correct: false },
        { text: "", image_url: "", is_correct: false },
        { text: "", image_url: "", is_correct: false },
        { text: "", image_url: "", is_correct: false },
      ],
    };
    setQuiz((prev) => ({ ...prev, questions: [...prev.questions, newQuestion] }));
  };

  const handleQuestionChange = (index, updatedQuestion) => {
    const updatedQuestions = [...quiz.questions];
    updatedQuestions[index] = updatedQuestion;
    setQuiz((prev) => ({ ...prev, questions: updatedQuestions }));
  };

  const handleSubmit = () => {
    fetch(`${API_ENDPOINTS.QUIZ}/`, {
      method: "POST",
      headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${sessionStorage.getItem("access_token")}`
       },
      body: JSON.stringify(quiz),
    })
      .then((res) => {
        if (!res.ok) throw new Error("Submission failed");
        return res.json();
      })
      .then((data) => {
        alert("Quiz submitted successfully!");
        navigate("/store");
      })
      .catch((err) => alert("Error: " + err.message));
  };

  return (
    <div className={styles["constructor-page"]}>
      <nav className={styles["nav-links"]}>
        <a href="/">Home</a>
        <a href="/join">Join</a>
        <a href="/create">Create</a>
      </nav>

      <div className={styles["constructor-content"]}>
        <div className={styles["constructor-header"]}>
          <input
            type="text"
            placeholder="Enter Quiz Title"
            value={quiz.title}
            onChange={(e) => setQuiz((prev) => ({ ...prev, title: e.target.value }))}
          />
          <span>{quiz.questions.length} questions</span>
        </div>

        {quiz.questions.map((question, index) => (
          <QuestionCard
            key={index}
            index={index}
            question={question}
            onChange={(updatedQuestion) => handleQuestionChange(index, updatedQuestion)}
          />
        ))}

        <div className={styles["add-question"]} onClick={handleAddQuestion}>＋</div>

        <button className={styles["submit-button"]} onClick={handleSubmit}>Submit</button>
      </div>
    </div>
  );
};

export default ConstructorPage