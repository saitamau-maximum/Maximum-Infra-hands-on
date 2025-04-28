import { GetAllRoomsResponse } from "../../types/GetAllRoomsResponse";

import styles from "./RoomList.module.css";

type RoomListProps = {
  rooms: GetAllRoomsResponse[];
};

export const RoomList = ({ rooms }: RoomListProps) => {
  return (
    <div className={styles.roomList}>
      <h2>ルーム一覧</h2>
      <ul>
        {rooms.map((room) => (
          <li key={room.id}>
            {room.name}
          </li>
        ))}
      </ul>
    </div>
  );
}