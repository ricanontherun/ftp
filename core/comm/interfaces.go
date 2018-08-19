package comm

type ClientInterface interface {
	Connect() error
	Close() error
}
