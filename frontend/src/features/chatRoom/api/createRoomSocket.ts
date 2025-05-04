import apiClient from "../../utils/apiClient";

export const createRoomSocket = (roomId: string): WebSocket | null => {
  try {
    return apiClient.websocket('/api/ws/' + roomId);
  } catch (err) {
    console.error('WebSocket作成失敗', err);
    return null;
  }
}