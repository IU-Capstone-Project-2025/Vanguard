import React, { useState } from "react";
import { QRCodeSVG } from "qrcode.react";
import { useParams, useNavigate } from "react-router-dom";
import styles from './styles/AskToJoinSession.module.css';

import { BASE_URL } from '../constants/api';

const AskToJoinSession = () => {
  const { sessionCode } = useParams();
  const [joinLink] = useState(`${BASE_URL}/wait/${sessionCode}`);
  const [copied, setCopied] = useState(false);
  const navigate = useNavigate();

  const handleCopyClick = async () => {
    try {
      await navigator.clipboard.writeText(joinLink);
      setCopied(true);
      setTimeout(() => {
        setCopied(false);
      }, 2000);
    } catch (err) {
      // console.error('Failed to copy text: ', err);
    }
  };

  const handlePlayClick = () => {
    navigate(`/sessionAdmin/${sessionCode}`);
  };

  return (
    <div className={styles["ask-join-container"]}>
      <div className={styles["ask-left"]}>
        <div className={styles["ask-title"]}>
          <h1>It's your code<br />for joining</h1>
          <div className={styles["join-code"]}>
            <span onClick={handleCopyClick}>
              #{sessionStorage.getItem("sessionCode")}
            </span>
            <span>
              {copied && <p>Link copied!</p>}
            </span>
          </div>
          <button className={styles["play-button"]} onClick={handlePlayClick}>▶ Play</button>
        </div>
      </div>
      <div className={styles["ask-right"]}>
        <div className={styles["qr-wrapper"]}>
          <QRCodeSVG value={joinLink} size={480} />
        </div>
      </div>
    </div>
  );
};

export default AskToJoinSession