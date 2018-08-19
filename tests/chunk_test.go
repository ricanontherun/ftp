package tests

import (
	"testing"
	"path/filepath"
	"ftp/core/reader"
	"path"
	"os"
	"math"
)

const (
	WordsFile string = "words"
	FirstFiveBytes string = "Lorem"
	NextSixBytes string = " ipsum"
)

var dataDir, _ = filepath.Abs("./data")

func TestNewChunkReader(t *testing.T) {
	// It should return an error if the file doesn't exist.
	_, err := reader.NewChunkReader(path.Join(dataDir, "idontexist"))
	if err == nil || !os.IsNotExist(err) {
		t.Fatalf("Expected to receive not exists error, instead received %s", err)
	}

	// It should not return an error when provided a valid file path.
	chunkReader, err := reader.NewChunkReader(path.Join(dataDir, WordsFile))
	defer chunkReader.Close()

	if err != nil {
		t.Fatalf("Expected nil error, received %s", err)
	}
}

func TestGetNextChunk(t *testing.T) {
	chunkReader := getWordsChunkReader(t)
	defer chunkReader.Close()

	// Read the first 5 byte chunk of the file ("Lorem")
	var readSize uint64 = 5
	chunk, _ := chunkReader.GetNextChunk(readSize)
	if chunk.Len != 5 {
		t.Fatalf("Received %d bytes, expecting %d", chunk.Len, readSize)
	}

	if string(chunk.Data) != FirstFiveBytes {
		t.Fatal("Failed to read 'there' from file")
	}

	// Reading the next 6 byte chunk should equal " ipsum".
	readSize = 6
	chunk, _ = chunkReader.GetNextChunk(readSize)

	if chunk.Len != 6 {
		t.Fatalf("Received %d bytes, expecting %d", chunk.Len, readSize)
	}

	if string(chunk.Data) != NextSixBytes {
		t.Fatal("Failed to read ' ipsum' from file")
	}

	chunkReader.Close()

	chunkReader = getWordsChunkReader(t) // We don't want to rely on ResetOffset
	defer chunkReader.Close()
	stat, _ := os.Stat(chunkReader.GetFilePath())

	// Read the entire file in one chunk.
	chunk, err := chunkReader.GetNextChunk(uint64(stat.Size()))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if int64(chunk.Len) != stat.Size() {
		t.Fatalf("Failed to read entire file, received %d bytes instead of %d", chunk.Len, stat.Size())
	}
}

func TestResetOffset(t *testing.T) {
	chunkReader := getWordsChunkReader(t)
	defer chunkReader.Close()

	if _, err := chunkReader.GetNextChunk(5); err != nil {
		t.Fatalf("Failed to read 5 bytes from file: %s", err)
	}

	if _, err := chunkReader.GetNextChunk(6); err != nil {
		t.Fatalf("Failed to read 6 bytes from file: %s", err)
	}

	if err := chunkReader.ResetOffset(); err != nil {
		t.Fatalf("Failed to ResetOffset(): %s", err)
	}

	// Considering we just reset the file offset to 0, the next 5 bytes should be lorem.
	if chunk, err := chunkReader.GetNextChunk(5); err != nil {
		t.Fatalf("Failed to read 6 bytes from file: %s", err)
	} else if string(chunk.Data) != "Lorem" {
		t.Fatalf("Failed to reset offset, next chunk yielded '%s'", chunk.Data)
	}
}

func TestRead(t *testing.T) {
	chunkReader := getWordsChunkReader(t)
	defer chunkReader.Close()

	stat, err := os.Stat(chunkReader.GetFilePath())
	if err != nil {
		t.Fatalf("Failed to stat file: %s", err)
	}

	var expectedChunks = math.Ceil(float64(stat.Size()) / float64(chunkReader.GetReadSize()))
	var actualChunks = 0.0

	err = chunkReader.Read(func(chunk *reader.Chunk) {
		actualChunks += 1
	})

	if err != nil {
		t.Fatalf("Recevied unexpected error: %e", err)
	}

	if actualChunks != expectedChunks {
		t.Fatalf("Failed to read entire file in chunks, got %f chunks, expected %f", actualChunks, expectedChunks)
	}
}

func getWordsChunkReader(t *testing.T) reader.ChunkReaderInterface {
	return getChunkReader(path.Join(dataDir, WordsFile), t)
}

func getChunkReader(filePath string, t *testing.T) reader.ChunkReaderInterface {
	chunkReader, err := reader.NewChunkReader(filePath)

	if err != nil {
		t.Fatalf("Failed to create chunkReader: %s", err)
	}

	return chunkReader
}