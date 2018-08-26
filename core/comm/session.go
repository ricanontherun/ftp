package comm

import (
	"ftp/core/file"
	"math"
	"os"
)

type SessionOptions struct {
	// The source path on the client.
	Source string

	// The destination path on the server.
	Destination string

	// Size of the file being transferred
	Size uint64

	// The number of chunks we intend to transfer.
	NumChunks uint64

	// The size (at most) of each chunk we intend to transfer.
	ChunkSize uint64
}

type SessionChunk struct {
	file.Chunk
	SequenceNum uint64
}

type SessionInterface interface {
}

type Session struct {
	Chunks []SessionChunk
}

func NewSessionOptions(transferOptions TransferOptions) (*SessionOptions, error) {
	stat, err := os.Stat(transferOptions.Source)

	if err != nil {
		return nil, err
	}

	sessionOptions := &SessionOptions{
		Source:      transferOptions.Source,
		Destination: transferOptions.Destination,
		Size:        uint64(stat.Size()),
		// TODO: Leverage the optimal read blk size here.
		ChunkSize: file.DefaultReadSize,
	}

	sessionOptions.NumChunks = uint64(math.Ceil(
		float64(sessionOptions.Size) / float64(sessionOptions.ChunkSize),
	))

	return sessionOptions, nil
}

func NewSession(sessionOptions *SessionOptions) (SessionInterface, error) {
	session := &Session{}

	session.Chunks = make([]SessionChunk, sessionOptions.NumChunks)

	return session, nil
}
