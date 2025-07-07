import React, { createContext, useContext, useRef, useCallback } from 'react';
import { API_ENDPOINTS } from '../constants/api';

const SessionWebSocketContext = createContext();

export const SessionWebSocketProvider = ({ children }) => {
  const wsRef = useRef(null);

  const connect = useCallback((token, onMessage) => {
    if (wsRef.current) return;

    const ws = new WebSocket(`${API_ENDPOINTS.SESSION_WS}?token=${token}`);
    ws.onmessage = onMessage;
    wsRef.current = ws;
  }, []);

  return (
    <SessionWebSocketContext.Provider value={{ wsRef, connect }}>
      {children}
    </SessionWebSocketContext.Provider>
  );
};

export const useSessionSocket = () => useContext(SessionWebSocketContext);
