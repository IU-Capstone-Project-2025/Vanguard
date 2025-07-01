import React, {useEffect, useRef, useState} from "react";
import { useNavigate } from "react-router-dom";
import "./styles/WaitGameStartAdmin.css";

const WaitGameStartAdmin = () => {
  const navigate = useNavigate();
  const sessionServiceWsRef = useRef(null);

  const [players, setPlayers] = useState([
    { id: 1, name: "Alice" },
    { id: 2, name: "Bob" },
    { id: 3, name: "Charlie" },
    { id: 4, name: "Diana" },
    { id: 5, name: "Eva" },
    { id: 6, name: "Frank" },
    { id: 7, name: "Grace" },
    { id: 8, name: "Henry" },
    { id: 9, name: "Isabella" },
  ]);

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
        // const players = JSON.parse(message.data);
        setPlayers(JSON.parse(message.data));
        console.log("📨 Received JSON message:",players);
      } catch (e){
        console.error("⚠️ Failed to parse incoming WebSocket message:", message.data);
      }
    }
  };

  useEffect(() => {
    connectToWebSocket(sessionStorage.getItem("jwt"))
  },[])
  useEffect(() => {},
      [players])

  const handleKick = (idToRemove) => {
    setPlayers(prev => prev.filter(player => player.id !== idToRemove));
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
          {players.map((player) => (
            <div key={player.id} className="player-box">
              <span>{player.name}</span>
              <button
                className="kick-button"
                onClick={() => handleKick(player.id)}
                title={`Kick ${player.name}`}
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
