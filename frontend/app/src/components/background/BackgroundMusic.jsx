import { useEffect, useRef } from "react";

const BackgroundMusic = () => {
  const audioRef = useRef(null);

  useEffect(() => {
    const audio = audioRef.current;
    audio.volume = 0.4;

    const tryPlay = () => {
      audio.play().catch(err => {
        console.log("Автовоспроизведение заблокировано:", err);
      });
    };

    document.addEventListener("click", tryPlay, { once: true });

    return () => {
      document.removeEventListener("click", tryPlay);
    };
  }, []);

  return (
    <audio ref={audioRef} loop>
      <source src="/audio/background-sound.mp3" type="audio/mpeg" />
    </audio>
  );
};

export default BackgroundMusic;
