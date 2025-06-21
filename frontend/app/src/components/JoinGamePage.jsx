import React from "react";
import './styles/styles.css'
import { useNavigate } from "react-router-dom";
import { useState } from "react";

const JoinGamePage = () => {
    const [code,setCode] = useState("")
    const navigate = useNavigate()

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
                        onChange={(e)=> setCode(e.target.value)}
                        className="code-input"
                    />
                    <div className="button-group">
                        <button id="play"
                                className="play-button"
                                onClick={
                                    (e) => {
                                        navigate('/play')
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