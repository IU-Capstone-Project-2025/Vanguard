// ShowLeaderBoardComponent.jsx
import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import './ShowLeaderBoardComponent.css';

const ShowLeaderBoardComponent = ({ leaderboardData, onClose }) => {
    const [players, setPlayers] = React.useState(
        sessionStorage.getItem('players') ? JSON.parse(sessionStorage.getItem('players')) : {}
    );

    return (
        <AnimatePresence>
        {leaderboardData && (
            <motion.div
            className="leaderboard-overlay"
            initial={{ y: '-100vh' }}
            animate={{ y: 0 }}
            exit={{ y: '-100vh' }}
            transition={{ type: 'spring', stiffness: 80 }}
            >
            <h1 className="leaderboard-title">Look! Here's our champions!</h1>
            {leaderboardData.users &&
            <div className="leaderboard-list">
                {leaderboardData.users.map((user, index) => (
                <div className="leaderboard-row" key={index}>
                    <span className="player-name">{players.find(player => player.id === user.user_id)?.name}</span>
                    <span className="player-score">{user.total_score}</span>
                </div>
                ))}
            </div>}
            <button className="leaderboard-next-button" onClick={onClose}>
                ▶ Next
            </button>
            </motion.div>
        )}
        </AnimatePresence>
    );
};

export default ShowLeaderBoardComponent;
