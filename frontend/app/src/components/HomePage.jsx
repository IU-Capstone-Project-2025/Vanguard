import React from "react";
import styles from './styles/HomePage.module.css';
import Homepage_image from './assets/homepage-image.png';
import { useNavigate } from "react-router-dom";

const HomePage = () => {
    const navigate = useNavigate();

    const handlePlayClick = () => {
        navigate('/join');
    };

    const handleCreateClick = () => {
        sessionStorage.removeItem('sessionCode');
        sessionStorage.removeItem('nickname');
        navigate('/create');
    };

    return (
        <div className={styles['homepage-main-content']}>
            <div className={styles['left-side']}>
                <div className={styles.title}>
                    <p>Just</p>
                    <h1>TRY.IT</h1>
                    <div className={styles['button-group']}>
                        <button
                            id="play"
                            className={styles['play-button']}
                            onClick={handlePlayClick}
                        >
                            <span>Play Game</span>
                        </button>
                        <button
                            id='create'
                            className={styles['create-button']}
                            onClick={handleCreateClick}
                        >
                            <span>Create Game</span>
                        </button>
                    </div>
                </div>
            </div>
            <div className={styles['right-side']}>
                <img src={Homepage_image} alt="Illustration of a boy" className={styles.image}/>
            </div>
        </div>
    );
};

export default HomePage;