import { useForm } from 'react-hook-form';
import { useRoomSocket } from '../hooks/useRoomSocket';
import { useParams } from 'react-router-dom';
import { ChatInput } from '../components';

import styles from './RoomPage.module.css';

type FormData = {
  message: string;
};

export const RoomPage = () => {
  const { roomId } = useParams<{ roomId: string }>();
  const { register, handleSubmit, reset } = useForm<FormData>();
  const { messages, sendMessage } = useRoomSocket(roomId ?? '');

  if (!roomId) {
    return <div>Room ID is required</div>;
  }

  const onSubmit = (data: FormData) => {
    if (data.message.trim()) {
      sendMessage(data.message);
      reset(); // フォームの値をリセット
    }
  };

  return (
    <div>
      <div>
        {messages.map((msg, index) => (
          <div key={index}>{msg}</div>
        ))}
      </div>
      <form onSubmit={handleSubmit(onSubmit)} className={styles.form}>
        <ChatInput.Field>
          <ChatInput.Input
            {...register('message')}
            placeholder="メッセージを入力"
          />
          <ChatInput.Button />
        </ChatInput.Field>
      </form>
    </div>
  );
};
