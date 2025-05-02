import { useState } from 'react';
import { useRoomSocket } from '../hooks/useRoomSocket';
import { useParams } from 'react-router-dom';

export const RoomPage = () => {
  const { roomId } = useParams<{ roomId: string }>();
  if (!roomId) {
    return <div>Room ID is required</div>;
  }
  const { messages, sendMessage } = useRoomSocket(roomId);
  const [input, setInput] = useState('');

  const handleSend = () => {
    if (input.trim()) {
      sendMessage(input);
      setInput('');
    }
  };

  return (
    <div>
      <div>
        {messages.map((msg, index) => (
          <div key={index}>
            {msg}
          </div>
        ))}
      </div>
      <input
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="メッセージを入力"
      />
      <button onClick={handleSend}>送信</button>
    </div>
  );
};
