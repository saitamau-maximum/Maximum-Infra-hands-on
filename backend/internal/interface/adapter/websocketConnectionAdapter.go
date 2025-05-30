package adapter

type ConnAdapter interface {
	ReadMessageFunc() (int, []byte, error)
	WriteMessageFunc(int, []byte) error
	CloseFunc() error
	ReadJSON(any) error
	WriteJSON(any) error
}
