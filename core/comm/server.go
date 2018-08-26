package comm

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

type FtpServer struct {
	sessions map[string]SessionInterface
}

func makeSessionToken(sessionOptions *SessionOptions) (string, error) {
	sha256 := sha256.New()

	if _, err := sha256.Write([]byte(sessionOptions.Destination)); err != nil {
		return "", err
	}

	return hex.EncodeToString(sha256.Sum(nil)), nil
}

func (server *FtpServer) Session(sessionOptions *SessionOptions, transferSession *TransferSession) error {
	var token string
	var tokenErr error

	if token, tokenErr = makeSessionToken(sessionOptions); tokenErr != nil {
		log.Println(tokenErr)
		return tokenErr
	}

	if _, exists := server.sessions[token]; exists {
		return ErrorDuplicateSession
	}

	transferSession.Token = token
	session, err := NewSession(sessionOptions)

	if err != nil {
		return ErrorSessionCreation
	}

	server.sessions[transferSession.Token] = session

	return nil
}

// Create a new server.
func NewServer() (*FtpServer, error) {
	server := &FtpServer{}

	server.sessions = make(map[string]SessionInterface)

	return server, nil
}
