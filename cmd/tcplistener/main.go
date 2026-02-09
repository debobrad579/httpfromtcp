package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		linesCh := getLinesChannel(conn)
		for line := range linesCh {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		currentLine := ""

		for {
			bytes := make([]byte, 8)

			if _, err := f.Read(bytes); err != nil {
				if err == io.EOF && currentLine != "" {
					linesCh <- currentLine
				}
				close(linesCh)
				f.Close()
				return
			}

			parts := strings.Split(string(bytes), "\n")

			for index, part := range parts {
				currentLine += part
				if index != len(parts)-1 {
					linesCh <- currentLine
					currentLine = ""
				}
			}
		}
	}()

	return linesCh
}
