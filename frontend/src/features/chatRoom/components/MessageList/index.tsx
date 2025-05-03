import { useCallback, useEffect, useRef, useState } from "react";
import { MessageResponse } from "../../type";
import { MessageListItem } from "./MessageListItem";

import styles from "./index.module.css";

interface MessageListProps {
  messages: MessageResponse[];
  fetchMore: () => void;
  hasNext: boolean;
}

export const MessageList = ({
  messages,
  fetchMore,
  hasNext,
}: MessageListProps) => {
  const scrollContainerRef = useRef<HTMLDivElement | null>(null);
  const bottomBoundaryRef = useRef<HTMLDivElement | null>(null);
  const [needFetchMore, setNeedFetchMore] = useState(false);

  const scrollObserver = useCallback(
    (node: Element) => {
      if (!scrollContainerRef.current) return;

      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting && hasNext) {
              setNeedFetchMore(true);
            }
          });
        },
        {
          root: scrollContainerRef.current,
          threshold: 1.0,
        }
      );

      observer.observe(node);
    },
    [hasNext]
  );

  useEffect(() => {
    if (bottomBoundaryRef.current) {
      scrollObserver(bottomBoundaryRef.current);
    }
  }, [scrollObserver]);

  useEffect(() => {
    if (needFetchMore) {
      fetchMore();
      setNeedFetchMore(false);
    }
  }, [needFetchMore, fetchMore]);

  return (
    <div
      ref={scrollContainerRef}
      className={styles.container}
    >
      {messages.map((message, index) => (
        <MessageListItem
          key={message.id}
          message={message}
          isOdd={index % 2 === 0}
        />
      ))}
      <div ref={bottomBoundaryRef} />
    </div>
  );
};
