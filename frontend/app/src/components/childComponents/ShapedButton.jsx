import React, { useEffect, useState } from 'react';
// import './ShapedButton.css';

const ShapedButton = ({ shape, text, onClick, className = '', style = {} }) => {
  const [pathData, setPathData] = useState(null);
  const [pathBBox, setPathBBox] = useState(null);
  const [clipId] = useState(`clip-${Math.random().toString(36).slice(2, 10)}`);

  useEffect(() => {
    const loadSvg = async () => {
      try {
        const res = await fetch(shape);
        const txt = await res.text();

        const parser = new DOMParser();
        const svgDoc = parser.parseFromString(txt, 'image/svg+xml');
        const path = svgDoc.querySelector('path');
        const svgElem = svgDoc.querySelector('svg');

        if (!path) {
          console.error('SVG does not contain a <path> element');
          return;
        }

        const d = path.getAttribute('d');
        const vb = svgElem?.getAttribute('viewBox')?.split(' ').map(Number) || [0, 0, 1000, 1000];

        setPathData(d);
        setPathBBox({ x: vb[0], y: vb[1], width: vb[2], height: vb[3] });
      } catch (err) {
        console.error('Failed to load or parse SVG', err);
      }
    };

    loadSvg();
  }, [shape]);

  if (!pathData || !pathBBox) return null;

  const scaleX = 1 / pathBBox.width;
  const scaleY = 1 / pathBBox.height;
  const translateX = -pathBBox.x * scaleX;
  const translateY = -pathBBox.y * scaleY;
  const transform = `translate(${translateX},${translateY}) scale(${scaleX},${scaleY})`;

  return (
    <div className={`shaped-button-wrapper ${className}`} style={style}>
      <svg className="hidden-svg" aria-hidden="true" focusable="false">
        <defs>
          <clipPath id={clipId} clipPathUnits="objectBoundingBox">
            <path d={pathData} transform={transform} />
          </clipPath>
        </defs>
      </svg>

      <button
        className="shaped-button"
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
