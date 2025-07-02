import React, { useEffect, useState ,useRef} from "react";
import { useNavigate } from "react-router-dom";
import Cookies from "js-cookie"

import './styles/styles.css';

const CreateSessionPage = () => {
  const navigate = useNavigate();

  const [selectedQuiz, setSelectedQuiz] = useState(null);
  const [quizzes, setQuizzes] = useState([]);
  const [search, setSearch] = useState("");
  const [sessionCode, setSessionCode] = useState(null);
  const realTimeWsRef = useRef(null);

  useEffect(() => {
    const fetchQuizzes = async () => {
      try {
        const response = await fetch("/api/quiz/");
        if (!response.ok) {
          throw new Error(`Network error: ${response.status}`);
        }
        const data = await response.json();
        if (!Array.isArray(data)) {
          throw new Error("Expected an array of quizzes.");
        }
        // Мапим массив, берем только id и title
        setQuizzes(data.map(quiz => ({ id: quiz.id, title: quiz.title })));
      } catch (error) {
        console.error("Error fetching quizzes:", error);
      }
    };
    fetchQuizzes();
  }, []);

   const handleQuizSelection = (quiz) => {
    sessionStorage.setItem('selectedQuizId', quiz.id);
    setSelectedQuiz(quiz);
  };

  // 🎯 POST-запрос к /api/session/sessions
  const createSession = async (quizId, userId) => {
    console.log("Creating session with quizId:", quizId, "and userId:", userId, "userName:", Cookies.get("user_nickname"));
    const response = await fetch("/api/session/sessions", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        "userId": userId,
        "userName": Cookies.get("user_nickname"),
        "quizId": quizId,
      }),
    });

    if (!response.ok) throw new Error("Failed to create session");

    const data = await response.json();
    return data; // возвращает объект вида: {"serverWsEndpoint": "string","jwt": "string", "sessionId":"string"}
  };

  // 🌐 Устанавливаем WebSocket-соединение с Real-time Service
  const connectToWebSocket = (token) => {
      realTimeWsRef.current = new WebSocket(`/api/session/ws?token=${token}`);
      realTimeWsRef.current.onopen = () => {
          console.log("✅ WebSocket connected with Session Service");
      };

  realTimeWsRef.current.onerror = (err) => {
    console.error("❌ WebSocket with Session Service error:", err);

  };
}


  const handlePlay = async () => {
    if (selectedQuiz) {
      const sessionData = await createSession(selectedQuiz.id, "AdminId");
      setSessionCode(sessionData.sessionId);
      // await connectToWebSocket(sessionData.jwt);

      sessionStorage.setItem('selectedQuizId', selectedQuiz.id);
      sessionStorage.setItem('sessionCode', sessionData.sessionId);
      sessionStorage.setItem('jwt', sessionData.jwt);

      navigate(`/ask-to-join/${sessionData.sessionId}`); // Navigate to the waiting page with the session code

    }
  };

  return (
    <div className="create-session-main-content">
      <div className="left-side">
        <div className="title">
          <h2>Now choose the quiz <br /> to start a game</h2>
          <div className="button-group">
            <button
              className="play-button"
              onClick={handlePlay}
              disabled={!selectedQuiz}
            >
              ▶ Play
            </button>
            <button
              className="enter-store-button"
              onClick={() => navigate("/store")}
            >
              + Enter quiz Store
            </button>
          </div>
        </div>
      </div>

      <div className="right-side">
        <div className="quiz-list-container">
          <div className="quiz-search-panel">
            <input
              type="text"
              placeholder="Search the quiz"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="quiz-search-input"
            />

            <div className="quiz-list">
              {quizzes
                .filter((quiz) =>
                  quiz.title.toLowerCase().includes(search.toLowerCase())
                )
                .map((quiz) => (
                  <div
                    key={quiz.id}
                    className={`quiz-item ${selectedQuiz?.id === quiz.id ? 'selected' : ''}`}
                    onClick={() => handleQuizSelection(quiz)}
                  >
                    {quiz.title}
                  </div>
                ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateSessionPage;
