import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from './styles/EnterNicknamePage.module.css';
import { API_ENDPOINTS } from "../constants/api.js";

const PlayGamePage = () => {
    const [nickname, setNickname] = useState("");
    const [error, setError] = useState(null);
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();
    const inputRef = useRef(null);

    useEffect(() => {
        inputRef.current?.focus();
    }, []);

    const joinSession = async (sessionCode, userName) => {
        try {
            setIsLoading(true);
            const response = await fetch(`${API_ENDPOINTS.SESSION}/join`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    "code": sessionCode,
                    "userName": userName
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || "Failed to join session");
            }

            return await response.json();
        } catch (error) {
            // console.error("Error joining session:", error);
            setError(error.message || "Failed to join session");
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

    const handlePlay = async () => {
        if (!nickname.trim()) {
            setError("Please enter a nickname to continue");
            return;
        }

        sessionStorage.setItem('nickname', nickname);
        
        try {
            const sessionCode = sessionStorage.getItem("sessionCode");
            if (sessionCode) {
                const sessionData = await joinSession(sessionCode, nickname);
                if (!sessionData?.sessionId) {
                    throw new Error("Invalid session data received");
                }
                sessionStorage.setItem('jwt', sessionData.jwt);
                navigate(`/wait/${sessionCode}`);
            } else {
                navigate('/join');
            }
        } catch (error) {
            // Error is already handled in joinSession
        }
    };

    return (
        <div className={styles['playgame-main-content']}>
            <div className={styles['left-side']}>
                <div className={styles.title}>
                    <h1>
                        <span>Who</span> are you today?
                    </h1>
                    <input 
                        type="text"
                        ref={inputRef}
                        placeholder="Enter your name here"
                        required
                        autoFocus
                        value={nickname}
                        onChange={(e) => {
                            const nick = e.target.value;
                            if (nick.length <= 16) {
                                setNickname(nick);
                                setError(null);
                            }
                        }}
                        className={styles['code-input']}
                        onKeyDown={(e) => e.key === "Enter" && handlePlay()}
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
        </div>
    );
};

export default PlayGamePage;