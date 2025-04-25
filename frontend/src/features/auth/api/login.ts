import { LoginFormData } from "../types/LoginFormDate";
import apiClient from "../../utils/apiClient";

export const Login = async (data: LoginFormData): Promise<void> => {
  const res = await apiClient.post(
    "/api/user/login",
    data,
  );

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "登録に失敗しました");
  }
};
