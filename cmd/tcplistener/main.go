package main

import (
	"fmt"
	"log"
	"net"

	"github.com/debobrad579/httpfromtcp/internal/http"
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

		req, err := http.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Print("Headers:\n")
		req.Headers.Range(func(fieldName, fieldValue string) bool {
			fmt.Printf("- %s: %s\n", fieldName, fieldValue)
			return true
		})
		fmt.Printf("Body:\n%s\n", string(req.Body))

		fmt.Println("Connection closed")
	}
}
