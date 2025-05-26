import apiClient from "../../utils/apiClient"

export const uploadImageApi = async (file: File): Promise<void> => {
  const formData = new FormData();
  formData.append("icon", file);

  const response = await apiClient.post(`/api/user/icon`, formData);
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || "画像のアップロードに失敗しました");
  }
};
