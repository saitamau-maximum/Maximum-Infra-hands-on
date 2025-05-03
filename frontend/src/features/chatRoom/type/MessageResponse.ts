export type MessageResponse = {
  id: string;
  user_id: string;
  sent_at: string; // RFC3339
  content: string;
};
