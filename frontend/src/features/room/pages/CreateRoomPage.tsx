import { CreateRoomForm } from "../components/CreateRoomForm/CreateRoomForm";

import styles from "./CreateRoomPage.module.css";

export const CreateRoomPage = () => {
  return (
    <div className={styles.container}>
      <h1>Create Room</h1>
      <CreateRoomForm />
    </div>
  );
}