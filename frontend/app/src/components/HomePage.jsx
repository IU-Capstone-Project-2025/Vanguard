import React from "react";
import './styles/styles.css'
import Homepage_image from './assets/homepage-image.png'
import { useNavigate } from "react-router-dom";

const HomePage = () => {
    const navigate = useNavigate()

    const handlePlayClick = () => {
        navigate('/join')
    }

    const handleCreateClick = () => {
        sessionStorage.removeItem('sessionCode'); // Clear any previous session code
        sessionStorage.removeItem('nickname'); // Clear any previous nickname
        navigate('/create')
    }

    return (
        <div className="homepage-main-content">
            <div className="left-side">
                <div className="title">
                    <h1>
                        Explore the InnoQuiz
                    </h1>
                    <div className="button-group">
                        <button id="play"
                                className="play-button"
                                onClick={handlePlayClick}
                            >
                            <span>Play Game</span>
                        </button>
                        <button id='create'
                                className="create-button"
                                onClick={handleCreateClick}
                            >
                            <span>Create Game</span>
                        </button>
                    </div>
                </div>
            </div>
            <div className="right-side">
                <img src={Homepage_image} alt="image with the boy"/>
            </div>
        </div>
    );

};

export default HomePage;