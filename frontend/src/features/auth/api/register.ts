import { RegisterFormData } from "../types/RegisterFormData";

export const Register = async (data: RegisterFormData): Promise<void> => {
  const res = await fetch("http://localhost:8080/api/user/register", {// TODO: 環境変数化
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
    credentials: "include",
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "登録に失敗しました");
  }
};
