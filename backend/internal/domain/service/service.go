// DI 層で使用するサービスの定義
package service

type Service struct {
	// IconStoreService はアイコンのストレージサービスです。
	IconStoreService IconStoreService

	// MessageCacheService はメッセージのキャッシュサービスです。
	MessageCacheService MessageCacheService

	// WebsocketManager はWebSocket接続を管理するマネージャーです。
	WebsocketManager WebsocketManager
}