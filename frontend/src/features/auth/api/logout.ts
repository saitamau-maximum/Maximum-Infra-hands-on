export const Logout = async (): Promise<void> => {
  const res = await fetch("http://localhost:8080/api/user/logout", {// TODO: 環境変数化
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });

  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || "ログアウトに失敗しました");
  }
};
