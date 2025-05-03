import { MessageResponse } from "../../type";

import styles from "./MessageListItem.module.css";

type MessageListItemProps = {
  message: MessageResponse;
  isOdd: boolean;
};

export const MessageListItem = ({ message, isOdd }: MessageListItemProps) => {
  const className = isOdd ? styles.odd : styles.even;

  return (
    <div className={className}>
      <p className={styles.message}>{message.content}</p>
    </div>
  );
};
