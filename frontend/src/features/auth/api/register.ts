import apiClient from "../../utils/apiClient";
import { RegisterFormData } from "../types/RegisterFormData";

export const Register = async (data: RegisterFormData): Promise<void> => {
  const res = await apiClient.post("/api/user/register", data);

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "登録に失敗しました");
  }
};
