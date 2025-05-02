export const parseRoomMessage = (event: MessageEvent): any => {
  try {
    return JSON.parse(event.data);
  } catch {
    return event.data; // 文字列のまま返す
  }
};
