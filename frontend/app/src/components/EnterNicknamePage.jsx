import React from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";

const PlayGamePage = () => {
    const [nickname,setNickname] = useState("")
    const navigate = useNavigate()


    const handlePlay = () => {
        if (nickname) {
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