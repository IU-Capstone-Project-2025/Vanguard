import React, {useEffect, useRef} from "react";
import './styles/styles.css'
import {useNavigate} from "react-router-dom";
import {useState} from "react";
import Cookies from "js-cookie";

const JoinGamePage = () => {
    const [code, setCode] = useState("");
    const navigate = useNavigate()
    const realTimeWsRef = useRef(null);

    // 🎯 POST-запрос к /api/session/join
    const joinSession = async (sessionCode, userId) => {
        console.log("Joining session with code:", sessionCode, "and userId:", userId);

        // Проверка на наличие кода с символом '#'
        const response = await fetch("/api/session/join", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({"code": sessionCode, "userId": userId, "userName":Cookies.get("user_nickname")}),
        });
        if (response.code === 400) {
            console.error("❌ Error joining session:", response.statusText);
            alert("Failed to join session. Please check the code and try again.");
            return;
        }
        const data = await response.json();
        return data; // возвращает объект вида: {"serverWsEndpoint": "string","jwt":"string", "sessionId":"string"}

    };
    useEffect(() => {
        console.log("code updated:", code);
    }, [code]);

    const handlePlay = async () => {
        if (code) {
            console.log("code updated:", code);
            const sessionData = await joinSession(code, "PlayerId")
            if (!sessionData || !sessionData.sessionId) {
                console.error("❌ Failed to join session or session ID is missing");
                alert("Failed to join session. Please check the code and try again.");
                return;
            }
            sessionStorage.setItem('sessionCode', code); // Store the session code in session storage
            sessionStorage.setItem('jwt', sessionData.jwt);
            navigate(`/wait/${code}`); // Navigate to the waiting page with the session code
        }
    };

    return (
        <div className="joingame-main-content">
            <div className="left-side">
                <div className="title">
                    <h1>
                        Ask your quiz creator for a code
                    </h1>
                    <input
                        type="text"
                        placeholder="enter a code here"
                        value={code}
                        onChange={(e) => {
                            const value = e.target.value.toUpperCase();
                            if (/^[A-Z0-9]*$/.test(value)){
                                setCode(value);
                        }}} // Remove '#' and convert to uppercase
                        required
                        autoFocus
                        pattern="^[A-Z0-9]+$" // Ensure only alphanumeric characters are
                        className="code-input"
                    />
                    <div className="button-group">
                        <button id="play"
                                className="play-button"
                                onClick={
                                    (e) => {
                                        handlePlay();
                                        e.preventDefault();
                                    }
                                }
                        >
                            <span>Play</span>
                        </button>
                    </div>
                </div>
            </div>
            <div className="right-side">

            </div>
        </div>
    )
};

export default JoinGamePage;