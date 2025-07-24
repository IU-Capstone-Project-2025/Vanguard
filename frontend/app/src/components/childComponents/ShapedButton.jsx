import React, { useEffect, useState } from 'react';
import styles from './ShapedButton.module.css';

const ShapedButton = ({ shape, text, onClick, className = '', style = {} }) => {
  const [svgData, setSvgData] = useState({ path: null, viewBox: null });
  const [clipId] = useState(`clip-${Math.random().toString(36).slice(2, 10)}`);

  useEffect(() => {
    const loadSvg = async () => {
      try {
        const response = await fetch(shape);
        const svgText = await response.text();
        
        const parser = new DOMParser();
        const svgDoc = parser.parseFromString(svgText, 'image/svg+xml');
        const path = svgDoc.querySelector('path');
        const svg = svgDoc.querySelector('svg');

        if (!path || !svg) {
          // console.error('SVG missing required elements');
          return;
        }

        setSvgData({
          path: path.getAttribute('d'),
          viewBox: svg.getAttribute('viewBox')?.split(' ').map(Number) || [0, 0, 1000, 1000]
        });
      } catch (error) {
        // console.error('Error loading SVG:', error);
      }
    };

    loadSvg();
  }, [shape]);

  if (!svgData.path || !svgData.viewBox) return null;

  const [x, y, width, height] = svgData.viewBox;
  const transform = `
    translate(${-x / width},${-y / height}) 
    scale(${1 / width},${1 / height})
  `;

  return (
    <div className={`${styles['shaped-button-wrapper']} ${className}`} style={style}>
      <svg className={styles['hidden-svg']} aria-hidden="true" focusable="false">
        <defs>
          <clipPath id={clipId} clipPathUnits="objectBoundingBox">
            <path d={svgData.path} transform={transform} />
          </clipPath>
        </defs>
      </svg>

      <button
        className={styles['shaped-button']}
        style={{
          clipPath: `url(#${clipId})`,
          WebkitClipPath: `url(#${clipId})`,
        }}
        onClick={onClick}
      >
        {text}
      </button>
    </div>
  );
};

export default ShapedButton;