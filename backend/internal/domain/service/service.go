// DI 層で使用するサービスの定義
package service

type Service struct {
	IconStoreService IconStoreService
	MessageCacheService MessageCacheService
	WebsocketManager WebsocketManager
}