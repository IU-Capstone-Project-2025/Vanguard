import React from 'react';
import styles from './FloatingBackground.module.css';
import Shape1 from './shapes/Shape1.svg';
import Shape2 from './shapes/Shape2.svg';
import Shape3 from './shapes/Shape3.svg';
import Shape4 from './shapes/Shape4.svg';
import Shape5 from './shapes/Shape5.svg';
import Shape6 from './shapes/Shape6.svg';
import Shape7 from './shapes/Shape7.svg';

const svgIcons = [Shape1, Shape2, Shape3, Shape4, Shape5, Shape6, Shape7];

const FloatingBackground = () => {
  const count = 12;
  const sectionWidth = 50 / count;
  const verticalRange = -400;
  const sectionHeight = verticalRange / count;

  return (
    <div className={styles['floating-background']}>
      {Array.from({ length: count }).map((_, i) => {
        const icon = svgIcons[i % svgIcons.length];
        const left = (i - count / 4) * sectionWidth + Math.random() * (sectionWidth * 0.6);
        const bottomBase = i * sectionHeight;
        const bottom = bottomBase + Math.random() * (sectionHeight * 0.6);
        const size = 480 + Math.random() * 50;
        const animationDelay = i * 2.5;
        
        return (
          <img
            key={i}
            src={icon}
            className={styles['floating-icon']}
            style={{
              left: `${left}%`,
              bottom: `${bottom}%`,
              width: `${size}px`,
              animationDelay: `${animationDelay}s`,
              zIndex: -1,
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