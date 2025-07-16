import React, { useEffect, useState } from 'react';
import './styles/FinalAdminPage.css';

const FinalAdminPage = () => {
  const [players, setPlayers] = useState([]);

  useEffect(() => {
    const storedPlayers = sessionStorage.getItem('players');
    if (!storedPlayers) return;

    try {
      const parsedPlayers = JSON.parse(storedPlayers);

      // Сортируем по убыванию score, затем назначаем позиции
      const sortedPlayers = [...parsedPlayers]
        .sort((a, b) => b.score - a.score)
        .slice(0, 3) // ← только трое
        .map((player, index) => ({ ...player, position: index + 1 }));


      setPlayers(sortedPlayers);
    } catch (err) {
      console.error('Ошибка при парсинге игроков:', err);
    }
    console.log('Final players:', players);
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
        <button className="primary-button" onClick={() => window.location.href = '/'}>
          Back to Home
        </button>
      </div>
    </div>
  );
};

export default FinalAdminPage;
