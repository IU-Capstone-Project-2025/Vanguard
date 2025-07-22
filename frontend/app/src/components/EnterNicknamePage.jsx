import React, {useEffect, useRef} from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";
// Import API_ENDPOINTS from its module (adjust the path as needed)
import { API_ENDPOINTS } from "../constants/api.js";

const PlayGamePage = () => {
    const [nickname,setNickname] = useState("")
    const navigate = useNavigate()
    const inputRef = useRef(null)

    useEffect(() => {
        inputRef.current.focus()
    }, []);
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
            if (response.statusCode === 400) {
                console.error("❌ Error joining session:", response.statusText);
                alert("Failed to join session. Please check the code and try again.");
                return;
            }else if (response.statusCode === 405) {
                console.error("❌ Error joining session:", response.statusText);
                return ;
            }else if (response.statusCode === 500) {
                console.error("❌ Error joining session:", response.statusText);
                return ;
            }
            const data = await response.json();
            return data; // возвращает объект вида: {"serverWsEndpoint": "string","jwt":"string", "sessionId":"string"}
    
        };

    const handlePlay = async () => {
        sessionStorage.setItem('nickname', nickname     )
        if (sessionStorage.getItem("sessionCode") !== undefined && nickname) {
            const code = sessionStorage.getItem("sessionCode");
            const sessionData = await joinSession(code, sessionStorage.getItem("nickname"))
            if (!sessionData || !sessionData.sessionId) {
                console.error("❌ Failed to join session or session ID is missing");
                alert("Failed to join session. Please check the code and try again.");
                return;
            }
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
            <div className="left-side">
                <div className="title">
                    <h1>
                        <span>Who</span> are you today?
                    </h1>
                    <input 
                        type="text"
                        ref={inputRef}
                        placeholder="enter the name here"
                        required
                        autoFocus
                        value={nickname}
                        onChange={(e)=> {
                            const nick = e.target.value;
                            if (nick.length<=16) {
                                setNickname(nick)
                        }}}
                        className="code-input"
                        onKeyDown={(e) => e.key === "Enter" && handlePlay()}
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
        </div>
    )
};

export default PlayGamePage;