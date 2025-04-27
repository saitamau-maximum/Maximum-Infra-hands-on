import apiClient from "../../utils/apiClient"

export const getAllRoomsApi = async () => {
  const response = await apiClient.get("/api/room");
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || "ルームの取得に失敗しました");
  }

  const data = await response.json();
  return data;
}