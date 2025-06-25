import React from 'react';
import { useNavigate } from 'react-router-dom';
import './styles/styles.css';

const NotFoundPage = () => {
    const navigate = useNavigate();

    return (
        <div className="notfound-page" style={{ textAlign: 'center', padding: '50px', backgroundColor: '#f0f0f0', minHeight: '100vh' }}>
            <div className='content-404' style={{ maxWidth: '600px', margin: '0 auto', backgroundColor: '#f8f9fa', padding: '20px', borderRadius: '10px' }}>
                <h1>404 - Page Not Found</h1>
                <p>It seems you've taken a wrong turn.</p>
                <p>Don't worry, we will add this page in further versions!</p>
                <button 
                    className="play-button" 
                    onClick={() => navigate('/')}
                >
                    Back to Home
                </button>
            </div>
        </div>
    );
};

export default NotFoundPage;
