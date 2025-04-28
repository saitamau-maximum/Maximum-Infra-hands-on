import apiClient from "../../utils/apiClient"
import { CreateRoomFormData } from "../types/CreateRoomFormDate"

export const CreateRoomApi = async (data : CreateRoomFormData) => {
  const response = await apiClient.post("/api/room", data);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || "ルームの作成に失敗しました");
  }
}