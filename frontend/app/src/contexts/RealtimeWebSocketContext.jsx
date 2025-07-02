import React, { createContext, useContext, useRef, useCallback } from 'react';

const RealtimeWebSocketContext = createContext();

export const RealtimeWebSocketProvider = ({ children }) => {
  const wsRef = useRef(null);

  const connect = useCallback((token, onMessage) => {
    if (wsRef.current) return;

    const ws = new WebSocket(`/api/realtime/ws?token=${token}`);
    ws.onmessage = onMessage;
    wsRef.current = ws;
  }, []);

  return (
    <RealtimeWebSocketContext.Provider value={{ wsRef, connect }}>
      {children}
    </RealtimeWebSocketContext.Provider>
  );
};

export const useRealtimeSocket = () => useContext(RealtimeWebSocketContext);
