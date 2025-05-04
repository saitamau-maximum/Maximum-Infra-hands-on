import { useEffect, useRef, useState, useCallback } from 'react';
import { createRoomSocket, parseRoomMessage, sendRoomMessage } from '../api';
import { MessageResponse } from '../type';

export const useRoomSocket = (roomId: string) => {
  const socketRef = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<MessageResponse[]>([]);

  useEffect(() => {
    const socket = createRoomSocket(roomId);
    if (!socket) {
      console.error('WebSocket creation failed');
      return;
    }
    socketRef.current = socket;

    socket.onmessage = (event) => {
      const data = parseRoomMessage(event);
      console.log('Received message:', data);
      // MessageResonse型に変換
      const msg: MessageResponse = {
        id: data.ID as string,
        user_id: data.UserID as string,
        sent_at: data.SentAt as string,
        content: data.Content as string,
      }
      setMessages((prev) => [...prev, msg]);
    };

    socket.onerror = (err) => console.error('WebSocket error:', err);
    socket.onclose = () => console.log('WebSocket closed');

    return () => {
      socket.close();
    };
  }, [roomId]);

  const sendMessage = useCallback((message: string) => {
    if (socketRef.current) {
      sendRoomMessage(socketRef.current, message);
    }
  }, []);

  return {
    messages,
    sendMessage,
    socket: socketRef.current,
  };
};
