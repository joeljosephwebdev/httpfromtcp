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
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Request line: \n - Method: %s\n - Target: %s\n - Version: %s\nHeaders:\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion)
		for key, value := range request.Headers {
			fmt.Printf(" - %s: %s\n", key, value)
		}
		fmt.Printf("Body: \n %s\n", request.Body)
		fmt.Println("connection closed")
	}
}
