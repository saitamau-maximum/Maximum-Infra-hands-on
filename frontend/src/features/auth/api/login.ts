import { LoginFormData } from "../types/LoginFormDate";
import { useNavigate } from "react-router-dom";

export const loginUser = async (data: LoginFormData): Promise<void> => {
  const navigate = useNavigate();
  const res = await fetch("http://localhost:8080/api/user/register", {
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
  navigate("/");
};
