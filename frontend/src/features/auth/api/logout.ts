import apiClient from "../../utils/apiClient";

export const Logout = async (): Promise<void> => {
  const res = await apiClient.post("/api/user/logout", null);

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "ログアウトに失敗しました");
  }
};
