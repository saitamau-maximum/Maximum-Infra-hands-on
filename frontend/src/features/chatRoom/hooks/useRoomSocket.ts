import { useEffect, useRef, useState, useCallback } from 'react';
import { createRoomSocket, parseRoomMessage, sendRoomMessage } from '../api';

export const useRoomSocket = (roomId: string) => {
  const socketRef = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<any[]>([]);

  useEffect(() => {
    const socket = createRoomSocket(roomId);
    if (!socket) {
      console.error('WebSocket creation failed');
      return;
    }
    socketRef.current = socket;

    socket.onmessage = (event) => {
      console.log('Received message:', event.data);
      const data = parseRoomMessage(event);
      setMessages((prev) => [...prev, data.Content]);
    };

    socket.onerror = (err) => console.error('WebSocket error:', err);
    socket.onclose = () => console.log('WebSocket closed');

    return () => {
      socket.close();
    };
  }, [roomId]);

  const sendMessage = useCallback((message: string) => {
    if (socketRef.current) {
      console.log('Sending message:', message);
      sendRoomMessage(socketRef.current, message);
    }
  }, []);

  return {
    messages,
    sendMessage,
    socket: socketRef.current,
  };
};
