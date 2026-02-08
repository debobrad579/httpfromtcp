package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		return
	}
	defer file.Close()

	linesCh := getLinesChannel(file)
	for line := range linesCh {
		fmt.Println("read: " + line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		currentLine := ""

		for {
			bytes := make([]byte, 8)

			if _, err := f.Read(bytes); err != nil {
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
