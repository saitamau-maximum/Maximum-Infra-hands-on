import { useForm } from "react-hook-form";
import { Form } from "../../../ui/Form";
import { CreateRoomFormData } from "../../types/CreateRoomFormDate";

import styles from "./CreateRoomForm.module.css";

export const CreateRoomForm = () => {
  const {
    register,
    handleSubmit,
  } = useForm<CreateRoomFormData>();
  const handleCreateRoom = async (data: CreateRoomFormData) => {
    try {
      // TODO: 部屋作成APIを呼び出す
      console.log("Room created:", data);
    } catch (error) {
      console.error("Room creation failed:", error);
    }
  };
  return (
    <form className={styles.form} onSubmit={handleSubmit(handleCreateRoom)}>
      <Form.Field>
        <Form.Label label="Room Name" />
        <Form.Input
          type="text"
          id="roomName"
          required
          placeholder="Enter room name"
          {...register("name")}
        />
        <Form.Button type="submit">
          Create Room
        </Form.Button>
      </Form.Field>
    </form>
  );
}