import React from "react";
import './styles/styles.css'
import Homepage_image from './assets/homepage-image.png'
import { useNavigate } from "react-router-dom";

const HomePage = () => {
    const navigate = useNavigate()

    const handlePlayClick = () => {
        navigate('/play')
    }

    const handleCreateClick = () => {
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
                            <span>Play</span>
                        </button>
                        <button id='create'
                                className="create-button"
                                onClick={handleCreateClick}
                            >
                            <span>Create</span>
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