import { useForm } from "react-hook-form";
import { Form } from "../../../ui/Form";
import { CreateRoomFormData } from "../../types/CreateRoomFormDate";

import styles from "./CreateRoomForm.module.css";
import { CreateRoomApi } from "../../api";
import { useNavigate } from "react-router-dom";

export const CreateRoomForm = () => {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
  } = useForm<CreateRoomFormData>();
  const handleCreateRoom = async (data: CreateRoomFormData) => {
    try {
      await CreateRoomApi(data);
      console.log("Room created:", data);
      navigate("/room");
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