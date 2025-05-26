const apiClient = {
  baseUrl: import.meta.env.VITE_API_BASE_URL, // ベースURLを指定

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
  post: async (endpoint: string, body: any, customHeaders?: HeadersInit) => {
    // FormData のときは Content-Type を指定しない（自動設定）
    const isFormData = body instanceof FormData;

    return apiClient.request(endpoint, {
      method: "POST",
      headers: isFormData
        ? undefined
        : customHeaders ?? { "Content-Type": "application/json" },
      body: isFormData ? body : JSON.stringify(body),
    });
  },

  // websocket接続
  websocket: (endpoint: string) => {
    const ws = new WebSocket(`${apiClient.baseUrl}${endpoint}`);
    return ws;
  }
};

export default apiClient;
