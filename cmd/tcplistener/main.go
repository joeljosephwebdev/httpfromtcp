package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joeljosephwebdev/httpfromtcp/internal/request"
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
		requestLine, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Request line: \n - Method: %s\n - Target: %s\n - Version: %s\n",
			requestLine.RequestLine.Method,
			requestLine.RequestLine.RequestTarget,
			requestLine.RequestLine.HttpVersion)
		fmt.Println("connection closed")
	}
}
