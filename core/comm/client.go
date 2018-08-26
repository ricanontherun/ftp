package comm

import (
	"fmt"
	"ftp/core/file"
	"log"
	"net/rpc"
)

type ftpClient struct {
	client *rpc.Client
}

func Connect(options ConnectionOptions) (ClientInterface, error) {
	client, err := rpc.DialHTTP("tcp", options.Target)

	if err != nil {
		return nil, err
	}

	return &ftpClient{client: client}, nil
}

func (client *ftpClient) MakeSession(sessionOptions *SessionOptions) (*TransferSession, error) {
	transferSession := new(TransferSession)

	err := client.client.Call(RPCMakeSession, sessionOptions, transferSession)
	if err != nil {
		log.Fatalln(err)
	}

	return transferSession, nil
}

func (client *ftpClient) Transfer(options TransferOptions) error {
	sessionOptions, sessionOptionsErr := NewSessionOptions(options)
	if sessionOptionsErr != nil {
		return sessionOptionsErr
	}

	transferSession, makeSessionErr := client.MakeSession(sessionOptions)
	if makeSessionErr != nil {
		return makeSessionErr
	}

	chunkReader, chunkReaderErr := file.NewChunkReader(options.Source)
	defer chunkReader.Close()

	if chunkReaderErr != nil {
		return chunkReaderErr
	}

	fmt.Println(transferSession.Token)

	chunkReader.Read(func(chunk *file.Chunk) {
		fmt.Println(string(chunk.Data))
	})

	return nil
}

func (client *ftpClient) Close() error {
	if client.client != nil {
		client.client.Close()
	}

	return nil
}
