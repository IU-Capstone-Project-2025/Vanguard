import React, {
  createContext,
  useContext,
  useRef,
  useCallback,
} from 'react';
import { API_ENDPOINTS } from '../constants/api';

const RealtimeWebSocketContext = createContext(null);

export const RealtimeWebSocketProvider = ({ children }) => {
  const wsRefRealtime = useRef(null);

  const connectRealtime = useCallback((token, onMessage = () => {}) => {
    if (
      wsRefRealtime.current &&
      (wsRefRealtime.current.readyState === WebSocket.OPEN ||
       wsRefRealtime.current.readyState === WebSocket.CONNECTING)
    ) {
      // Already connected — update message handler
      wsRefRealtime.current.onmessage = onMessage;
      return;
    }

    const ws = new WebSocket(
      `${API_ENDPOINTS.REALTIME_WS}?token=${token}`
    );

    // ws.onopen = () => // console.log('✅ Realtime WebSocket connected');
    ws.onmessage = onMessage;
    // ws.onerror = (e) => // console.error('⚠️ Realtime WebSocket error:', e);
    ws.onclose = () => {
      // console.warn('ℹ️ Realtime WebSocket closed');
    };

    wsRefRealtime.current = ws;
  }, []);

  return (
    <RealtimeWebSocketContext.Provider value={{ wsRefRealtime, connectRealtime, closeWsRefRealtime: () => {
      if (wsRefRealtime.current) {
        wsRefRealtime.current.close(
          1000, // Normal closure
          'Closing Realtime WebSocket connection'
        );
      }
    }}}>
      {children}
    </RealtimeWebSocketContext.Provider>
  );
};

export const useRealtimeSocket = () => {
  const context = useContext(RealtimeWebSocketContext);
  if (!context) {
    throw new Error(
      'useRealtimeSocket must be used within a RealtimeWebSocketProvider'
    );
  }
  return context;
};
