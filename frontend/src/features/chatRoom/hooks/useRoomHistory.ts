import { useEffect, useState } from "react";
import { getMessageHistory } from "../api/getMessageHistory";
import { MessageResponse } from "../type";

// 1ページあたりのメッセージ数
const MESSAGE_LIMIT = 20;

export const useRoomHistory = (roomId: string, initialLoad = true) => {
  const [messages, setMessages] = useState<MessageResponse[]>([]);
  const [hasNext, setHasNext] = useState(true);
  const [loading, setLoading] = useState(false);
  const [nextBeforeSentAt, setNextBeforeSentAt] = useState<string | null>(null);

  const loadMore = async () => {
    if (loading || !hasNext) return;

    setLoading(true);
    try {
      console.log("履歴取得開始", roomId, nextBeforeSentAt);
      const res = await getMessageHistory({
        roomId,
        limit: MESSAGE_LIMIT,
        beforeSentAt: nextBeforeSentAt ?? undefined,
      });
      console.log("履歴取得", res);

      setMessages((prev) => [...prev, ...res.messages]);
      console.log("履歴取得後", res.messages);
      setNextBeforeSentAt(res.next_before_sent_at);
      console.log("次の取得時刻", res.next_before_sent_at);
      setHasNext(res.has_next);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (initialLoad) {
      loadMore();
    }
    // 初回読み込みのため loadMore ではなく nextBeforeSentAt に依存させない
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialLoad]);

  return {
    messages,
    loadMore,
    loading,
    hasNext,
  };
};
