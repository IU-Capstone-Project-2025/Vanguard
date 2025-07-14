import React, { useEffect, useState } from 'react';
import './ShapedButton.css';

const ShapedButton = ({ shape, text, onClick }) => {
  const [pathData, setPathData] = useState(null);
  const [pathBBox, setPathBBox] = useState(null);
  const clipId = `clip-${btoa(shape).replace(/[^a-zA-Z0-9]/g, '')}`;

  useEffect(() => {
    const loadSvg = async () => {
      try {
        const res = await fetch(shape);
        const txt = await res.text();

        const parser = new DOMParser();
        const svgDoc = parser.parseFromString(txt, 'image/svg+xml');

        const path = svgDoc.querySelector('path');
        if (!path) {
          console.error('SVG does not contain a <path> element');
          return;
        }

        const d = path.getAttribute('d');

        // Получаем viewBox из SVG, если есть
        const svgElem = svgDoc.querySelector('svg');
        let viewBox = svgElem?.getAttribute('viewBox');
        let vb = null;

        if (viewBox) {
          // Превращаем строку "minX minY width height" в массив чисел
          vb = viewBox.split(' ').map(Number);
        } else {
          // Если viewBox нет — пытаемся вычислить bounding box вручную (не в DOM, а в JS)
          // Тут нет прямого способа без рендера, можно сделать грубое приближение
          // Но лучше требовать viewBox в SVG
          console.warn('SVG has no viewBox, scaling may be inaccurate');
          vb = [0, 0, 1000, 1000]; // примерный fallback
        }

        setPathData(d);
        setPathBBox({ x: vb[0], y: vb[1], width: vb[2], height: vb[3] });
      } catch (err) {
        console.error('Failed to load or parse SVG', err);
      }
    };

    loadSvg();
  }, [shape]);

  if (!pathData || !pathBBox) return null;

  // Считаем масштаб для objectBoundingBox (0..1) — нужно трансформировать из viewBox в (0,0,1,1)
  const scaleX = 1 / pathBBox.width;
  const scaleY = 1 / pathBBox.height;
  const translateX = -pathBBox.x * scaleX;
  const translateY = -pathBBox.y * scaleY;

  // transform: сначала сдвигаем, чтобы minX,minY стали 0,0, потом масштабируем по ширине и высоте к 1
  const transform = `translate(${translateX},${translateY}) scale(${scaleX},${scaleY})`;

  return (
    <div className="shaped-button-wrapper">
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
