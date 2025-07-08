import React, {
  createContext,
  useContext,
  useRef,
  useCallback,
} from 'react';
import { API_ENDPOINTS } from '../constants/api';

const SessionWebSocketContext = createContext(null);

export const SessionWebSocketProvider = ({ children }) => {
  const wsRefSession = useRef(null);

  const connectSession = useCallback((token, onMessage = () => {}) => {
    if (
      wsRefSession.current &&
      (wsRefSession.current.readyState === WebSocket.OPEN ||
       wsRefSession.current.readyState === WebSocket.CONNECTING)
    ) {
      // Уже подключён — просто переустанавливаем обработчик
      wsRefSession.current.onmessage = onMessage;
      return;
    }

    const ws = new WebSocket(
      `${API_ENDPOINTS.SESSION_WS}?token=${token}`
    );

    ws.onopen = () => console.log('✅ Session WebSocket connected');
    ws.onmessage = onMessage;
    ws.onerror = (e) => console.error('⚠️ Session WebSocket error:', e);
    ws.onclose = () => {
      console.warn('ℹ️ Session WebSocket closed');
    };

    wsRefSession.current = ws;
  }, []);

  return (
    <SessionWebSocketContext.Provider value={{ wsRefSession, connectSession, closeWsRefSession: () => {
      if (wsRefSession.current) {
        wsRefSession.current.close(
          1000, // Normal closure
          'Closing Session WebSocket connection'
        );
      }
    }}}>
      {children}
    </SessionWebSocketContext.Provider>
  );
};

export const useSessionSocket = () => {
  const context = useContext(SessionWebSocketContext);
  if (!context) {
    throw new Error(
      'useSessionSocket must be used within a SessionWebSocketProvider'
    );
  }
  return context;
};
