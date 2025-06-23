// コネクションを隠蔽するためのアダプター
// 具体実装はbackend/internal/infrastructure/adapterImpl/connAdapterImpl
package adapter

type ConnAdapter interface {
	// ReadMessageFunc はメッセージ([]byte)を読み取る関数です。
	ReadMessageFunc() (int, []byte, error)

	// WriteMessageFunc はメッセージ([]byte)を書き込む関数です。
	WriteMessageFunc(int, []byte) error

	// CloseFunc は接続を閉じる関数です。
	CloseFunc() error

	// ReadJSON は JSON データを読み取る関数です。
	ReadJSON(any) error

	// WriteJSON は JSON データを書き込む関数です。
	WriteJSON(any) error
}
