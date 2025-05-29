// ログ出力のためのアダプターインターフェース
// 具体実装は/infrastructure/adapterImpl/loggerAdapterImpl
package adapter 

type LoggerAdapter interface {
	// Info は情報レベルのログを出力します。
	Info(msg string, args ...any)

	// Debug はデバッグレベルのログを出力します。
	Warn(msg string, args ...any)

	// Warn は警告レベルのログを出力します。
	Error(msg string, args ...any)
}