import { Link } from "react-router-dom";
import { GetAllRoomsResponse } from "../../types/GetAllRoomsResponse"

import styles from "./RoomListItem.module.css";

type RoomListItemProps = {
  room: GetAllRoomsResponse;
};

export const RoomListItem = ({ room }: RoomListItemProps) => {
  return (
    <Link to={`/room/${room.id}`} className={styles.roomListItem}>
      <div className={styles.name}>{room.name}</div>
    </Link>
  )
}
