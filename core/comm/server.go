package comm

type FtpServer struct{}

func (server *FtpServer) Connect(num int, ret *int) error {
	*ret = num + 1
	return nil
}
