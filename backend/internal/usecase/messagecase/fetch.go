package messagecase

import (
	"context"
	"time"

	"example.com/infrahandson/internal/domain/entity"
)

type GetMessageHistoryInRoomRequest struct {
	RoomID       entity.RoomID
	Limit        int
	BeforeSentAt time.Time
}

type GetMessageHistoryInRoomResponse struct {
	Messages         []*entity.Message
	NextBeforeSentAt time.Time
	HasNext          bool
}

func (uc *MessageUseCase) GetMessageHistoryInRoom(ctx context.Context, req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error) {
	// まずはキャッシュからの取得を試みる
	messages, err := uc.msgCache.GetRecentMessages(ctx, req.RoomID)
	if err != nil {
		return GetMessageHistoryInRoomResponse{}, err
	}

	// キャッシュが空の場合、キャッシュ利用をスキップしてDBへ
	if len(messages) > 0 {
		// キャッシュの中で最も新しい・古いメッセージの時刻を調べる
		earliest := messages[0].GetSentAt()
		latest := messages[0].GetSentAt()
		for _, t := range messages[1:] {
			temp := t.GetSentAt()
			if temp.Before(earliest) {
				earliest = temp
			}
			if temp.After(latest) {
				latest = temp
			}
		}

		// キャッシュが使えるか判定（キャッシュの最新が BeforeSentAt よりも古い）
		if latest.Before(req.BeforeSentAt) {
			return GetMessageHistoryInRoomResponse{
				Messages:         messages,
				NextBeforeSentAt: earliest,
				HasNext:          len(messages) >= req.Limit,
			}, nil
		}
	}

	// メッセージ履歴を取得（キャッシュが使えない場合）
	messages, nextBeforeSentAt, hasNext, err := uc.msgRepo.GetMessageHistoryInRoom(
		ctx,
		req.RoomID,
		req.Limit,
		req.BeforeSentAt,
	)
	if err != nil {
		return GetMessageHistoryInRoomResponse{}, err
	}

	return GetMessageHistoryInRoomResponse{
		Messages:         messages,
		NextBeforeSentAt: nextBeforeSentAt,
		HasNext:          hasNext,
	}, nil
}
