package main

import (
	"errors"
	"ftp/core/comm"
	"log"
	"os"
)

var (
	errorMissingCommand      = errors.New("missing command")
	errorMissingTargetServer = errors.New("missing target server")
	errorMissingSourcePath   = errors.New("missing source path")
	errorInvalidSourcePath   = errors.New("invalid source path")
)

const (
	commandSend = "send"
)

func validateSendArgs(args []string, argMap map[string]string) error {
	if len(args) == 0 {
		return errorMissingSourcePath
	}

	argMap["sourcePath"] = args[0]

	if _, err := os.Stat(argMap["sourcePath"]); os.IsNotExist(err) {
		return errorInvalidSourcePath
	}

	if len(args) == 2 {
		argMap["destinationPath"] = args[1]
	}

	return nil
}

func parseClientArguments() (map[string]string, error) {
	args := os.Args[1:]

	if len(args) == 0 {
		return nil, errorMissingCommand
	}

	// Must have at least a command and the target
	if len(args) < 2 {
		return nil, errorMissingTargetServer
	}

	command := args[0]
	commandArgs := args[2:]

	argMap := make(map[string]string)

	var err error
	switch command {
	case commandSend:
		err = validateSendArgs(commandArgs, argMap)
		break
	}

	if err != nil {
		return nil, err
	}

	argMap["command"] = command
	argMap["target"] = args[1]

	return argMap, nil
}

func performSend(args map[string]string) error {
	// Collect some metadata on the file transfer. Size, chunks, chunkSize etc.
	// This will be used by the server to establish a connection session.

	// Connect to the server.
	client, err := comm.Connect(comm.ConnectionOptions{
		Target: args["target"],
	})

	if err != nil {
		return err
	}

	// Transfer the file.
	return client.Transfer(comm.TransferOptions{
		Source:      args["sourcePath"],
		Destination: args["destinationPath"],
	})
}

func main() {
	args, parseError := parseClientArguments()

	if parseError != nil {
		log.Fatalf("Failed to parse command: %s\n", parseError)
	}

	if sendErr := performSend(args); sendErr != nil {
		log.Fatalln(sendErr)
	}
}
