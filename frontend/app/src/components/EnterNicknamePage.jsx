import React from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";

const PlayGamePage = () => {
    const [nickname,setNickname] = useState("")
    const navigate = useNavigate()

    return (
        <div className="playgame-main-content">
            
                <div className="title">
                    <h1>
                        Now enter your nickname
                    </h1>
                    <input 
                        type="text" 
                        placeholder="enter the name here"
                        value={nickname}
                        onChange={(e)=> setNickname(e.target.value)}
                        className="code-input"
                    />
                    <div className="button-group">
                        <button id="play"
                                className="play-button"
                                onClick={
                                    (e) => {
                                        navigate('/wait}')
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