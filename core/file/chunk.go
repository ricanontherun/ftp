package file

import (
	"io"
	"os"
	"syscall"
)

const (
	DefaultReadSize uint64 = 1024
)

type ChunkReaderOpts struct {
	ReadSize uint64
}

type ChunkReaderInterface interface {
	GetNextChunk(uint64) (*Chunk, error)
	Read(func(*Chunk)) error

	GetReadSize() uint64
	SetReadSize(uint64)
	GetFilePath() string
	ResetOffset() error
	Close()
}

type Chunk struct {
	Data []byte
	Len  int
}

type chunkReader struct {
	filePath   string
	fileHandle *os.File
	fileStat   *syscall.Stat_t
	readSize   uint64
	opts       ChunkReaderOpts
}

func NewChunkReader(filePath string) (ChunkReaderInterface, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	reader := &chunkReader{
		filePath:   filePath,
		fileHandle: file,
		readSize:   DefaultReadSize,
	}

	return reader, nil
}

func (reader *chunkReader) Read(f func(*Chunk)) error {
	var chunk *Chunk = nil
	var readError error = nil

	if resetErr := reader.ResetOffset(); resetErr != nil {
		return resetErr
	}

	for {
		if chunk, readError = reader.GetNextChunk(reader.GetReadSize()); readError != nil {
			break
		}

		f(chunk)
	}

	if readError == io.EOF {
		return nil
	}

	return readError
}

func (reader *chunkReader) GetNextChunk(readBytes uint64) (*Chunk, error) {
	chunk := &Chunk{
		Data: make([]byte, readBytes),
		Len:  0,
	}

	bytesRead, err := reader.fileHandle.Read(chunk.Data)

	if err != nil {
		return nil, err
	}

	chunk.Len = bytesRead
	return chunk, nil
}

func (reader *chunkReader) Close() {
	if reader.fileHandle != nil {
		reader.fileHandle.Close()
	}
}

func (reader *chunkReader) ResetOffset() error {
	_, err := reader.fileHandle.Seek(0, io.SeekStart)

	return err
}

func (reader *chunkReader) GetReadSize() uint64 {
	return reader.readSize
}

func (reader *chunkReader) SetReadSize(readSize uint64) {
	reader.readSize = readSize
}

func (reader *chunkReader) GetFilePath() string {
	return reader.filePath
}
