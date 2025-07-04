import React, {useEffect, useRef, useState} from "react";
import { useNavigate } from "react-router-dom";
import "./styles/WaitGameStartAdmin.css";

const WaitGameStartAdmin = () => {
  const navigate = useNavigate();
  const sessionServiceWsRef = useRef(null);

  const [players, setPlayers] = useState(new Map());

  // 🌐 Устанавливаем WebSocket-соединение с Session Service
  const connectToWebSocket = (token) => {
    let serverWsEndpoint = "ws://localhost:8081/ws";
    sessionServiceWsRef.current = new WebSocket(`${serverWsEndpoint}?token=${token}`);
    sessionServiceWsRef.current.onopen = () => {
      console.log("✅ WebSocket connected with Session Service");
    };

    sessionServiceWsRef.current.onerror = (err) => {
      console.error("❌ WebSocket with Session Service error:", err);

    };

    // получение сообщения от session service

    sessionServiceWsRef.current.onmessage = (message) => {
      try {
        const incoming = JSON.parse(message.data); // { user123: "Alice", user456: "Bob" }

        setPlayers((prevMap) => {
          const updatedMap = new Map(prevMap); // Копируем старую Map

          for (const [userId, name] of Object.entries(incoming)) {
            if (!updatedMap.has(userId)) {
              updatedMap.set(userId, name); // Добавляем только новых
            }
          }

          return updatedMap;
        });

        console.log("📨 Received JSON message:",incoming);
      } catch (e){
        console.error("⚠️ Failed to parse incoming WebSocket message:", message.data);
      }
    }
  };

  useEffect(() => {
    connectToWebSocket(sessionStorage.getItem("jwt"))
  },[])

  const handleKick = (idToRemove) => {
    setPlayers(prev =>{
      const updated = new Map(prev);
      updated.delete(idToRemove);
      return updated;
    });
    // TODO: отправить на backend сигнал о кике игрока по id
  };

  const toNextQuestion = async (sessionCode) => {

    const response = await fetch(`/api/session/session/${sessionCode}/nextQuestion`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ "code": sessionCode}),
    });

    if (!response.ok) throw new Error("Failed to get to the next question session");

  };

  const handleStart = async () => {
    const sessionCode = sessionStorage.getItem('sessionCode');
    await toNextQuestion(sessionCode)
    navigate(`/game-controller/${sessionCode}`);
  };

  const handleTerminate = () => {
    navigate("/");
  };

  return (
      <div className="wait-admin-container">
        <div className="wait-admin-panel">
          <h1>Now let's wait your friends</h1>
          <div className="admin-button-group">
            <button onClick={handleStart}>▶ Start</button>
            <button onClick={handleTerminate}>▶ Terminate</button>
          </div>
          <div className="players-grid">
            {Array.from(players.entries()).map(([id, name]) => (
                <div key={id} className="player-box">
                  <span>{name}</span>
                  <button
                      className="kick-button"
                      onClick={() => handleKick(id)}
                      title={`Kick ${name}`}
                  >
                    ❌
                  </button>
                </div>
            ))}
          </div>
        </div>
      </div>
  );
};

export default WaitGameStartAdmin;
