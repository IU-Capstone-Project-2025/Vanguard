const isDev = process.env.NODE_ENV === 'development';

export const BASE_URL = isDev
  ? process.env.REACT_APP_BASE_URL || 'http://localhost:3000'
  : window.location.origin;

export const API_ENDPOINTS = {
  AUTH: isDev
    ? process.env.REACT_APP_AUTH_API || 'http://localhost:8000/api'
    : '/api/auth',
  QUIZ: isDev
    ? process.env.REACT_APP_QUIZ_API || 'http://localhost:8001/api'
    : '/api/quiz',
  SESSION: isDev
    ? process.env.REACT_APP_SESSION_API || 'http://localhost:8081'
    : '/api/session',
  LEADERBOARD: isDev
    ? process.env.REACT_APP_LEADERBOARD_API || 'http://localhost:8082'
    : '/api/leaderboard',
  SESSION_WS: isDev
    ? process.env.REACT_APP_SESSION_WS || 'ws://localhost:8081/ws'
    : '/api/session/ws',
  REALTIME_WS: isDev
    ? process.env.REACT_APP_REALTIME_WS || 'ws://localhost:8080/ws'
    : '/api/realtime/ws',
};
