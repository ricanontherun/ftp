package comm

type ClientInterface interface {
	MakeSession(*SessionOptions) (*TransferSession, error)
	Transfer(TransferOptions) error
	Close() error
}
