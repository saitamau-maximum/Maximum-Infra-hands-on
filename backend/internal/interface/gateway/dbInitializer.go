// データベースの初期化インターフェース
// 具体実装はbackend/internal/infrastructure/gatewayImpl/db
package gateway

import "github.com/jmoiron/sqlx"

type DBInitializer interface {
	// DBInitializer はデータベースの初期化インターフェースです。
	Init() (*sqlx.DB, error)

	// InitSchema はデータベーススキーマを初期化・マイグレーションします。
	InitSchema(db *sqlx.DB) error
}
