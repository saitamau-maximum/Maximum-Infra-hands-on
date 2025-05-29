// パスワードのハッシュ化と照合を行う機能のラッパー
// 具体実装は/infrastructure/adapterImpl/hasherAdapterImpl
package adapter

type HasherAdapter interface {
	// HashPassword はパスワードをハッシュ化します。
	HashPassword(password string) (string, error)

	// ComparePassword はハッシュ化されたパスワードと入力されたパスワードを比較します。
	ComparePassword(hashedPassword, password string) (bool, error)
}