import { useState, useEffect } from 'react';
import { getAllRoomsApi } from '../api';
import { GetAllRoomsResponse } from '../types/GetAllRoomsResponse';

export const useRooms = () => {
  const [rooms, setRooms] = useState<GetAllRoomsResponse[]>([]); // ルームのデータを格納
  const [loading, setLoading] = useState<boolean>(true); // ローディング状態
  const [error, setError] = useState<string | null>(null); // エラーメッセージ

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        setLoading(true);
        setError(null); // エラーメッセージをリセット
        const roomsData = await getAllRoomsApi(); // API呼び出し

        const rooms: GetAllRoomsResponse[] = roomsData.map((room: any) => ({
          id: room.room_id,
          name: room.name,
        })); // 型を指定
        
        setRooms(rooms); // 取得したデータを設定
      } catch (err: any) {
        setError(err.message || "ルームの取得に失敗しました");
      } finally {
        setLoading(false);
      }
    };

    fetchRooms();
  }, []); // コンポーネントのマウント時に1度だけ実行

  return { rooms, loading, error };
};
