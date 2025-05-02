CREATE INDEX idx_messages_room_id_sent_at ON messages(room_id, sent_at DESC);
