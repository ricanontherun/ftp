package comm

import "errors"

// String constants - mostly RPC names.
const (
	RPCMakeSession string = "FtpServer.Session"
)

var (
	ErrorDuplicateSession error = errors.New("Duplicate session")
	ErrorSessionCreation  error = errors.New("Failed to create session")
)
