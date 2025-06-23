// DI 層で扱うほかの層を組み立てるために必要な構造体を定義する
package adapter

type Adapter struct {
	// HasherAdapter はハッシュ化を行うアダプターです。
	HasherAdapter HasherAdapter

	// TokenServiceAdapter はトークンの生成や検証を行うアダプターです。
	TokenServiceAdapter TokenServiceAdapter

	// LoggerAdapter はロギングを行うアダプターです。
	LoggerAdapter LoggerAdapter

	// Upgrader はWebSocketのアップグレードを行うアダプターです。
	Upgrader WebSocketUpgraderAdapter
}