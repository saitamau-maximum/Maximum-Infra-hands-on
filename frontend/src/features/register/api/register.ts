import { RegisterFormData } from "../types/RegisterFormData";
// TODO: 送った後のトークン処理など
export const registerUser = async (data: RegisterFormData): Promise<void> => {
  console.log(data);
  const res = await fetch("http://localhost:8080/api/user/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "登録に失敗しました");
  }
};
