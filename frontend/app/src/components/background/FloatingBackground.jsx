import React from 'react';
import './FloatingBackground.css';

import Shape1 from './shapes/Shape1.svg'
import Shape2 from './shapes/Shape2.svg'
import Shape3 from './shapes/Shape3.svg'
import Shape4 from './shapes/Shape4.svg'
import Shape5 from './shapes/Shape5.svg'
import Shape6 from './shapes/Shape6.svg'
import Shape7 from './shapes/Shape7.svg'


const svgIcons = [Shape1, Shape2, Shape3, Shape4, Shape5, Shape6, Shape7];

const FloatingBackground = () => {
  return (
    <div className="floating-background">
      {Array.from({ length: 7 }).map((_, i) => {
        const icon = svgIcons[i % svgIcons.length];
        const left = Math.random() * 100;
        const animationDelay = Math.random() * 70;
        const size = 600 + Math.random() * 40;

        // Пропустить центр (40%–60%)

        return (
          <img
            key={i}
            src={icon}
            className="floating-icon"
            style={{
              left: `${left}%`,
              width: `${size}px`,
              animationDelay: `${animationDelay}s`,
              transform: 'rotate(0deg)'
            }}
            alt=""
          />
        );
      })}
    </div>
  );
};

export default FloatingBackground;
