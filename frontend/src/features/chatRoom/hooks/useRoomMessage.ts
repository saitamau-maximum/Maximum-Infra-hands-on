import { useMemo } from 'react';
import { useRoomHistory } from './useRoomHistory';
import { useRoomSocket } from './useRoomSocket';
import { MessageResponse } from '../type';

// メッセージ配列を sent_at の降順（新しい順）でソート
const mergeMessages = (history: MessageResponse[], realtime: MessageResponse[]) => {
  const all = [...history, ...realtime];
  const uniqueMap = new Map<string, MessageResponse>();

  for (const msg of all) {
    uniqueMap.set(msg.id, msg); // 同じIDのメッセージは上書き（重複排除）
  }

  return Array.from(uniqueMap.values()).sort((a, b) => {
    return new Date(b.sent_at).getTime() - new Date(a.sent_at).getTime(); // 降順
  });
};

export const useRoomMessages = (roomId: string) => {
  const {
    messages: historyMessages,
    loadMore,
    loading,
    hasNext,
  } = useRoomHistory(roomId);

  const {
    messages: realtimeMessages,
    sendMessage,
    socket,
  } = useRoomSocket(roomId);

  // 履歴とリアルタイムメッセージを統合
  const messages = useMemo(() => {
    return mergeMessages(historyMessages, realtimeMessages);
  }, [historyMessages, realtimeMessages]);

  return {
    messages,        // 統合済み
    loadMore,        // 無限スクロール用
    loading,         // 読み込み中
    hasNext,         // さらに履歴があるか
    sendMessage,     // WebSocket送信
    socket,          // WebSocket本体
  };
};
