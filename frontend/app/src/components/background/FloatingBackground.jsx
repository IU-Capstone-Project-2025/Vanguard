import React from 'react';
import './FloatingBackground.css';

const svgIcons = [
  'Alien.svg',
  'Arrow.svg',
  'Clover.svg',
  'Cookie6.svg',
  'Corona.svg',
  'Ghosty.svg',
  'Pentagon.svg'
];

const FloatingBackground = () => {
  return (
    <div className="floating-background">
      {Array.from({ length: 25 }).map((_, i) => {
        const icon = svgIcons[i % svgIcons.length];
        const left = Math.random() * 100;
        const animationDelay = Math.random() * 10;
        const size = 20 + Math.random() * 40;

        // Пропускаем центральную область (40% - 60%)
        if (left > 40 && left < 60) return null;

        return (
          <img
            key={i}
            src={`/images/${icon}`}
            className="floating-icon"
            style={{
              left: `${left}%`,
              width: `${size}px`,
              animationDelay: `${animationDelay}s`,
              transform: `rotate(${30}deg)`
            }}
            alt=""
          />
        );
      })}
    </div>
  );
};

export default FloatingBackground;
