package main

import (
	"ftp/core/comm"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	server := &comm.FtpServer{}
	rpc.Register(server)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":33344")
	if err != nil {
		log.Fatalln(err)
	}

	http.Serve(listener, nil)
}
