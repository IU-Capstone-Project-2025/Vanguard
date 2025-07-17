import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './styles/FinalAdminPage.css';
// Import API_ENDPOINTS from its module (adjust the path as needed)
import { API_ENDPOINTS } from '../constants/api';

const FinalAdminPage = () => {
  const [players, setPlayers] = useState([]);
  const [sessionCode, setSessionCode] = useState(sessionStorage.getItem('sessionCode') || '');
  const navigate = useNavigate();

  const endSession = async (code) => {
    try {
      const response = await fetch(`${API_ENDPOINTS.SESSION}/session/${code}/end`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) throw new Error(`Failed to end session with code: ${code}`);

      sessionStorage.removeItem('sessionCode');
      sessionStorage.removeItem('quizData');
      sessionStorage.removeItem('currentQuestion');
      sessionStorage.removeItem('players');
      sessionStorage.removeItem('leaders');
      
      navigate('/');
    } catch (error) {
      console.error('Error ending the session:', error);
    }
  }

  useEffect(() => {
    const storedLeaders = sessionStorage.getItem('leaders');
    const storedPlayers = sessionStorage.getItem('players');
    if (!storedLeaders || !storedPlayers) return;

    try {
      const leaders = JSON.parse(storedLeaders);
      const players = JSON.parse(storedPlayers);

      // Создаём map: user_id → name
      const nameMap = {};
      players.forEach(p => {
        if (p.name !== 'Admin') {
          nameMap[p.id] = p.name;
        }
      });

      // Соединяем лидеров с именами
      const processed = leaders
        .filter(l => nameMap[l.user_id]) // пропускаем тех, кого нет в nameMap (напр., Admin)
        .sort((a, b) => b.score - a.score)
        .slice(0, 3)
        .map((leader, index) => ({
          ...leader,
          name: nameMap[leader.user_id],
          position: index + 1,
        }));

      setPlayers(processed);
      console.log('Final leaders with names:', processed);
    } catch (err) {
      console.error('Error parsing leaders or players:', err);
    }
  }, []);




  return (
    <div className="final-page-container">
      <h1 className="final-title">Meet the Champions!</h1>
      <div className="podium">
        {players
          .sort((a, b) => a.position - b.position)
          .map((player, idx) => (
            <div key={idx} className={`podium-block position-${player.position}`}>
              <span className="player-name">{player.name}</span>
            </div>
        ))}
      </div>
      <div className='buttons-container'>
        <button className="final-primary-button " onClick={() => endSession(sessionCode)}>
          Back to Home
        </button>
      </div>
    </div>
  );
};

export default FinalAdminPage;
