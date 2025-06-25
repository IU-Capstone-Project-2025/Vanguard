import React from "react";
import { useState } from "react";
import {QRCodeSVG} from "qrcode.react";
import { useParams } from "react-router-dom";
import "./styles/AskToJoinSession.css"; // we'll style separately
import { useNavigate } from "react-router-dom";

const AskToJoinSession = () => {
  // Mocked backend data
  const { sessionCode } = useParams();
  const [joinLink, setJoinLink] = React.useState(`https://localhost:3000/join/${sessionCode}`);
  const [copied, setCopied] = useState(false);
  const navigate = useNavigate();



  const handleCopyClick = async () => {
        try {
            await navigator.clipboard.writeText(joinLink);
            setCopied(true);
            setTimeout(() => {
                setCopied(false);
            }, 2000); // Hide message after 2 seconds
        } catch (err) {
            console.error('Failed to copy text: ', err);
        }
    };

  const handlePlayClick = () => {
      navigate(`/sessionAdmin/${sessionCode}`);
  };

  return (
    <div className="ask-join-container">
      <div className="ask-left">
        <div className="ask-title">
            <h1>It's your code<br />for joining</h1>
            <div className="join-code">
                <span onClick={handleCopyClick}>
                    #{sessionCode}
                </span>
                <span>
                    {copied && <p>Link copied!</p>}
                </span>
                </div>
            <button className="play-button" onClick={handlePlayClick}>▶ Play</button>
        </div>
      </div>
      <div className="ask-right">
        <div className="qr-wrapper">
          <QRCodeSVG value={joinLink} size={480} />
        </div>
      </div>
    </div>
  );
};

export default AskToJoinSession;
