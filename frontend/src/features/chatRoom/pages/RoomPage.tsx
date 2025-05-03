import { useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { ChatInput } from '../components';

import styles from './RoomPage.module.css';
import { useRoomMessages } from '../hooks/useRoomMessage';
import { MessageList } from '../components/MessageList';

type FormData = {
  message: string;
};

export const RoomPage = () => {
  const { roomId } = useParams<{ roomId: string }>();
  if (!roomId) return <div>Room ID is required</div>;

  const { register, handleSubmit, reset } = useForm<FormData>();
  const { messages, loadMore, sendMessage, hasNext } = useRoomMessages(roomId);
  console.log('messages', messages);

  const onSubmit = (data: FormData) => {
    if (data.message.trim()) {
      sendMessage(data.message);
      reset(); // フォームの値をリセット
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit(onSubmit)} className={styles.form}>
        <ChatInput.Field>
          <ChatInput.Input
            {...register('message')}
            placeholder="メッセージを入力"
          />
          <ChatInput.Button />
        </ChatInput.Field>
      </form>
      <MessageList
        messages={messages}
        fetchMore={loadMore}
        hasNext={hasNext}
      />
    </div>
  );
};
