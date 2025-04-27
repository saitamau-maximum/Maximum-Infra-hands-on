import { Link } from "react-router-dom";
import { RoomList } from "../components";
import { useRooms } from "../hooks/useRooms";

import styles from "./RoomListPage.module.css";

export const RoomListPage = () => {
  const { rooms, loading, error } = useRooms();
  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error}</div>;
  }
  return (
    <div className={styles.container}>
      <h1>Room List</h1>
      <Link to="/room/create" className={styles.createRoomLink}>
        Create Room
      </Link >
      <RoomList rooms={rooms} />
    </div>
  );
}