import React from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './styles/NotFoundPage.module.css';

const NotFoundPage = () => {
    const navigate = useNavigate();

    return (
        <div className={styles['notfound-page']}>
            <div className={styles['content-404']}>
                <h1>404 - Page Not Found</h1>
                <p>It seems you've taken a wrong turn.</p>
                <p>Don't worry, we will add this page in further versions!</p>
                <button 
                    className={styles['home-button']} 
                    onClick={() => navigate('/')}
                >
                    Back to Home
                </button>
            </div>
        </div>
    );
};

export default NotFoundPage;