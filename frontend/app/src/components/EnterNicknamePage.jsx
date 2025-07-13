import React from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";
// Import API_ENDPOINTS from its module (adjust the path as needed)
import { API_ENDPOINTS } from "../constants/api.js";

const PlayGamePage = () => {
    const [nickname,setNickname] = useState("")
    const navigate = useNavigate()

    const joinSession = async (sessionCode, userName) => {
            console.log("Joining session with code:", sessionCode, "and username:", userName);
    
            // Проверка на наличие кода с символом '#'
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
            if (response.code === 400) {
                console.error("❌ Error joining session:", response.statusText);
                alert("Failed to join session. Please check the code and try again.");
                return;
            }
            const data = await response.json();
            return data; // возвращает объект вида: {"serverWsEndpoint": "string","jwt":"string", "sessionId":"string"}
    
        };

    const handlePlay = async () => {
        sessionStorage.setItem('nickname', nickname     )
        if (sessionStorage.getItem("sessionCode") !== undefined && nickname) {
            const code = await sessionStorage.getItem("sessionCode");
            const sessionData = await joinSession(code, sessionStorage.getItem("nickname"))
            if (!sessionData || !sessionData.sessionId) {
                console.error("❌ Failed to join session or session ID is missing");
                alert("Failed to join session. Please check the code and try again.");
                return;
            }
            sessionStorage.setItem('sessionCode', code); // Store the session code in session storage
            sessionStorage.setItem('jwt', sessionData.jwt);
            navigate(`/wait/${code}`);
        }else if (nickname && !sessionStorage.getItem("sessionCode")) {
            sessionStorage.setItem('nickname', nickname);
            navigate('/join');
        }
        else {
            alert("Please enter a nickname to continue.");
        }
    }
    return (
        <div className="playgame-main-content">
            
                <div className="title">
                    <h1>
                        Now enter your nickname
                    </h1>
                    <input 
                        type="text" 
                        placeholder="enter the name here"
                        required
                        autoFocus
                        value={nickname}
                        onChange={(e)=> setNickname(e.target.value)}
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
    )
};

export default PlayGamePage;