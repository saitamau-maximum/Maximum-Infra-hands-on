// DI 層で扱うほかの層を組み立てるために必要な構造体を定義する
package adapter

type Adapter struct {
	HasherAdapter HasherAdapter
	TokenServiceAdapter TokenServiceAdapter
	LoggerAdapter LoggerAdapter
	Upgrader WebSocketUpgraderAdapter
}