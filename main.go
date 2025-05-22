package main

import (
	"fmt"
	"strings"

	"github.com/joeljosephwebdev/httpfromtcp/internal/request"
)

func main() {
	r, err := request.RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	if err != nil {
		return
	}
	fmt.Printf("{\nMethod: %s\n", r.RequestLine.Method)
	fmt.Printf("Version: %s\n", r.RequestLine.HttpVersion)
	fmt.Printf("Target: %s\n}\n", r.RequestLine.RequestTarget)

	r, err = request.RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	if err != nil {
		return
	}

	fmt.Printf("{\nMethod: %s\n", r.RequestLine.Method)
	fmt.Printf("Version: %s\n", r.RequestLine.HttpVersion)
	fmt.Printf("Target: %s\n}\n", r.RequestLine.RequestTarget)
}
