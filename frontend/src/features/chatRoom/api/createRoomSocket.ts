export const createRoomSocket = (roomId: string): WebSocket | null => {
  // TODO: 環境変数
  try {
    return new WebSocket('ws://localhost:8080/api/ws/' + roomId);
  } catch (err) {
    console.error('WebSocket作成失敗', err);
    return null;
  }
}