import apiClient from "../../utils/apiClient";
import { MessageResponse } from "../type";

type GetMessageHistoryParams = {
  roomId: string;
  limit?: number;
  beforeSentAt?: string; // RFC3339形式
};


type GetMessageHistoryResponse = {
  messages: MessageResponse[];
  next_before_sent_at: string;
  has_next: boolean;
};

export const getMessageHistory = async ({
  roomId,
  limit,
  beforeSentAt,
}: GetMessageHistoryParams) => {
  try {
    const response = await apiClient.get(
      `/api/message/${roomId}?limit=${limit}&before_sent_at=${beforeSentAt}`
    )

    const data = await response.json();
    // response型に合わせる
    const formattedData: GetMessageHistoryResponse = {
      messages: data.messages.map((message: MessageResponse) => ({
        id: message.id,
        user_id: message.user_id,
        sent_at: message.sent_at,
        content: message.content,
      })),
      next_before_sent_at: data.next_before_sent_at,
      has_next: data.has_next,
    };
    return formattedData;
  }
  catch (error) {
    console.error("Error fetching message history:", error);
    throw new Error("メッセージ履歴の取得に失敗しました");
  }
}

