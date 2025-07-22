import React, { useEffect, useRef, useState } from "react";
import styles from './styles/JoinGamePage.module.css';
import { useNavigate } from "react-router-dom";
import { API_ENDPOINTS } from '../constants/api';
import nicknameIcon from './assets/nickname-page.svg';

const JoinGamePage = () => {
    const [code, setCode] = useState("");
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();
    const inputRef = useRef(null);

    const joinSession = async (sessionCode) => {
        try {
            setIsLoading(true);
            const response = await fetch(`${API_ENDPOINTS.SESSION}/validate`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    "code": sessionCode
                }),
            });

            if (!response.ok) {
                throw new Error("Invalid session code");
            }
            return true;
        } catch (error) {
            console.error("Error joining session:", error);
            setError("Invalid session code. Please check and try again.");
            return false;
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        inputRef.current?.focus();
    }, []);

    const handlePlay = async () => {
        if (!code.trim()) {
            setError("Please enter a session code");
            return;
        }

        const isValid = await joinSession(code);
        if (isValid) {
            sessionStorage.setItem('sessionCode', code);
            navigate('/enter-nickname');
        }
    };

    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            handlePlay();
        }
    };

    return (
        <div className={styles['joingame-main-content']}>
            <div className={styles['left-side']}>
                <div className={styles.title}>
                    <h1>
                        Got a <span className={styles.code}>code</span>?
                        <br/>
                        Time to jump in!
                    </h1>
                    <input
                        type="text"
                        ref={inputRef}
                        placeholder="Enter a code here"
                        value={code}
                        onChange={(e) => {
                            const value = e.target.value.toUpperCase();
                            if (/^[A-Z0-9]*$/.test(value)) {
                                setCode(value);
                                setError(null);
                            }
                        }}
                        onKeyDown={handleKeyDown}
                        required
                        autoFocus
                        pattern="^[A-Z0-9]+$"
                        className={styles['code-input']}
                    />
                    {error && <div className={styles.error}>{error}</div>}
                    <div className={styles['button-group']}>
                        <button
                            id="play"
                            className={styles['play-button']}
                            onClick={handlePlay}
                            disabled={isLoading}
                        >
                            {isLoading ? "Joining..." : <span>Play</span>}
                        </button>
                    </div>
                </div>
            </div>
            <div className={styles['right-side']}>
                <div className={styles['right-side-content']}>
                    <img 
                        src={nicknameIcon} 
                        alt="Nickname icon" 
                        className={styles['nickname-icon']}
                    />
                </div>
            </div>
        </div>
    );
};

export default JoinGamePage;