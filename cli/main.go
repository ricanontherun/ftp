package main

import (
	"errors"
	"fmt"
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

func main() {
	args, err := parseClientArguments()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(args)
	return

	client, err := comm.Connect(comm.ConnectionOptions{
		Host: "localhost",
		Port: 33344,
	})

	if err != nil {
		log.Fatalln(err)
	}

	client.Connect()
}
