package comm

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
)

type ftpClient struct {
	client *rpc.Client
}

func Connect(options ConnectionOptions) (ClientInterface, error) {
	client, err := rpc.DialHTTP("tcp", options.Host+":"+strconv.Itoa(options.Port))

	if err != nil {
		return nil, err
	}

	return &ftpClient{client: client}, nil
}

func (client *ftpClient) Connect() error {
	var reply int

	err := client.client.Call("FtpServer.Connect", 1200, &reply)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(reply)
	return nil
}

func (client *ftpClient) Close() error {
	if client.client != nil {
		client.client.Close()
	}

	return nil
}
