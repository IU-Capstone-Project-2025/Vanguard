import React, {useEffect, useRef} from "react";
import './styles/styles.css'
import {useNavigate} from "react-router-dom";
import {useState} from "react";
import Cookies from "js-cookie";
import { API_ENDPOINTS } from '../constants/api';
import nicknameIcon from './assets/nickname-page.svg'; // Import the nickname icon if needed

const JoinGamePage = () => {
    const [code, setCode] = useState("");
    const navigate = useNavigate()
    const realTimeWsRef = useRef(null);
    const inputRef = useRef(null);

    // 🎯 POST-запрос к /api/session/validate
    const joinSession = async (sessionCode) => {
        console.log("Joining session with code:", sessionCode);

        // Проверка на наличие кода с символом '#'
        const response = await fetch(`${API_ENDPOINTS.SESSION}/validate`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "code": sessionCode
            }),
        });

        if (response.ok) {
            return true
        }else if (response.statusCode === 400) {
            console.error("❌ Error joining session:", response.statusText);
            alert("Failed to join session. Please check the code and try again.");
        }else if (response.statusCode === 405) {
            console.error("❌ Error joining session:", response.statusText);
        }else if (response.statusCode === 500) {
            console.error("❌ Error joining session:", response.statusText);
        }
        return false;



    };
    useEffect(() => {
        inputRef.current.focus();
    }, []);

    useEffect(() => {
        console.log("code updated:", code);
    }, [code]);

    const handlePlay = async () => {
        if (code) {
            console.log("code updated:", code);
            const sessionData = await joinSession(code)
            if (!sessionData) {
                console.error("❌ Failed to join session or session ID is missing");
                alert("Failed to join session. Please check the code and try again.");
                return;
            }
            sessionStorage.setItem('sessionCode', code); // Store the session code in session storage
            navigate(`/enter-nickname`); // Navigate to the waiting page with the session code
        }
    };

    return (
        <div className="joingame-main-content">
            <div className="left-side">
                <div className="title">
                    <h1>
                        Got a <span className="code">code</span>?
                        <br/>
                        Time to jump in!
                    </h1>
                    <input
                        type="text"
                        ref={inputRef}
                        placeholder="enter a code here"
                        value={code}
                        onChange={(e) => {
                            const value = e.target.value.toUpperCase();
                            if (/^[A-Z0-9]*$/.test(value)){
                                setCode(value);
                        }}} // Remove '#' and convert to uppercase
                        onKeyDown={(e) => e.key === "Enter" && handlePlay()}
                        required
                        autoFocus
                        pattern="^[A-Z0-9]+$" // Ensure only alphanumeric characters are
                        className="code-input"
                    />
                    <div className="button-group">
                        <button id="play"
                                className="play-button"
                                onClick={handlePlay}
                        >
                            <span>Play</span>
                        </button>
                    </div>
                </div>
            </div>
            <div className="right-side">
                <div className="right-side-content">
                    <img src={nicknameIcon} alt="nickname icon" className="nickname-icon" sizes="(max-width: 600px) 70vw, 30vw"/>
                </div>
            </div>
        </div>
    )
};

export default JoinGamePage;