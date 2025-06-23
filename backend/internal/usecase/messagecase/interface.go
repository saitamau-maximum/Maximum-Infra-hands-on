package messagecase

import "context"

type MessageUseCaseInterface interface {
	// GetMessageHisotyroInRoom: 一定数のメッセージ履歴を取得する
	GetMessageHistoryInRoom(ctx context.Context, req GetMessageHistoryInRoomRequest) (GetMessageHistoryInRoomResponse, error)
}

// Note: メッセージ作成はWebsocketのUseCase内で行われる（チャット通信との同期）
// そのため、ここには含めない
