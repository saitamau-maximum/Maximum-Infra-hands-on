import { Icon } from "../../../ui/icon";
import apiClient from "../../../utils/apiClient";
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
      <div className={styles.item}>
        <Icon src={`${apiClient.baseUrl}/api/user/icon/${message.user_id}`} size={24}/>
        <p className={styles.message}>{message.content}</p>
      </div>
    </div>
  );
};
