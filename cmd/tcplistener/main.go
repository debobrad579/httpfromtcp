package main

import (
	"fmt"
	"log"
	"net"

	"github.com/debobrad579/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Could not set up listener:", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection failed")
			continue
		}
		fmt.Println("Connection established")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		fmt.Println("Connection closed")
	}
}
