import React, { createContext, useContext, useRef, useCallback } from 'react';

const SessionWebSocketContext = createContext();

export const SessionWebSocketProvider = ({ children }) => {
  const wsRef = useRef(null);

  const connect = useCallback((token, onMessage) => {
    if (wsRef.current) return;

    const ws = new WebSocket(`/api/session/ws?token=${token}`);
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
