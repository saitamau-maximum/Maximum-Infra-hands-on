export const sendRoomMessage = (
  socket: WebSocket,
  message: string
): void => {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({
      content: message,
    }));
  } else {
    console.warn('WebSocket is not open. Message not sent:', message);
  }
};
