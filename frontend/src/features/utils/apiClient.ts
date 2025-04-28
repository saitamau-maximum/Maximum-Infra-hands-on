const apiClient = {
  baseUrl: `http://localhost:8080`, // ベースURLを指定

  // 基本的なリクエスト処理
  request: async (endpoint: string, options: RequestInit = {}) => {
    const res = await fetch(`${apiClient.baseUrl}${endpoint}`, {
      ...options,
      credentials: "include", // Cookieを送信
    });

    if (!res.ok) {
      const error = await res.json();
      throw new Error(error.message || "APIリクエストに失敗しました");
    }

    return res
  },

  // GETリクエスト専用
  get: async (endpoint: string) => {
    return apiClient.request(endpoint, { method: "GET" });
  },

  // POSTリクエスト専用
  post: async (endpoint: string, body: any) => {
    console.log(body);
    if (body == null) {
      return apiClient.request(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      credentials: "include",
      // Cookieを送信
      });
    }
    
    return apiClient.request(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
  },
};

export default apiClient;
