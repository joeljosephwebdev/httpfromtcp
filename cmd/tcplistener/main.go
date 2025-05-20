package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Listening on port 42069...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		fmt.Println("connection accepted")
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		currentLineContents := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLineContents + parts[i]
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
