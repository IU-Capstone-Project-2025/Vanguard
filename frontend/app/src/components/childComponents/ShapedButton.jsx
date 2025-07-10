import React, { useEffect, useState } from 'react';
import './ShapedButton.css';

const ShapedButton = ({ shape, text, onClick }) => {
  const [pathElement, setPathElement] = useState(null);
  const clipId = `clip-${btoa(shape).replace(/[^a-zA-Z0-9]/g, '')}`; // уникальный ID из пути

  useEffect(() => {
    const loadSvg = async () => {
      try {
        const res = await fetch(shape);
        console.log(res)
        const txt = await res.text();

        const parser = new DOMParser();
        const svgDoc = parser.parseFromString(txt, 'image/svg+xml');
        const path = svgDoc.querySelector('path');
        // console.log(txt)

        if (!path) {
          console.error('SVG does not contain a <path> element');
          return;
        }

        // Преобразуем DOM path в JSX path
        const pathJSX = (
          <path
            d={path.getAttribute('d')}
            transform="scale(0.003378378, 0.003164557)" // подгон под objectBoundingBox
          />
        );

        setPathElement(pathJSX);
      } catch (err) {
        console.error('Failed to load or parse SVG', err);
      }
    };

    loadSvg();
  }, [shape]);

  if (!pathElement) return null; // Не рендерим кнопку, пока SVG не загружен

  return (
    <div className="shaped-button-wrapper">
      <svg className="hidden-svg">
        <defs>
          <clipPath id={clipId} clipPathUnits="objectBoundingBox">
            {pathElement}
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
