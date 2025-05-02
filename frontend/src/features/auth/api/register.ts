import apiClient from "../../utils/apiClient";
import { RegisterFormData } from "../types/RegisterFormData";

type RegisterParams = {
  data: RegisterFormData;
  refetch: () => void;
};

export const Register = async ({data, refetch}: RegisterParams): Promise<void> => {
  const res = await apiClient.post("/api/user/register", data);
  if (res.ok) {
    refetch();
  }
  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "登録に失敗しました");
  }
};
