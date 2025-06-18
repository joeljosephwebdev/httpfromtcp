package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joeljosephwebdev/httpfromtcp/internal/headers"
	"github.com/joeljosephwebdev/httpfromtcp/internal/request"
	"github.com/joeljosephwebdev/httpfromtcp/internal/response"
	"github.com/joeljosephwebdev/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		writeBadRequest(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		writeServerError(w, req)
		return
	}

	headers := headers.NewHeaders()
	w.WriteStatusLine(response.StatusCodeSuccess)
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	body := response.BuildResponseBody(response.StatusCodeSuccess, "Your request was an absolute banger.")
	w.WriteBody(body)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	req_url := fmt.Sprintf("http://httpbin.org/%s", path)
	headers := headers.NewHeaders()

	resp, err := http.Get(req_url)
	if err != nil {
		writeServerError(w, req)
		log.Printf("Error making request to httpbin: %v\n", err)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCode(resp.StatusCode))

	// copy headers from resp to headers
	for k, v := range resp.Header {
		for _, vv := range v {
			headers.Add(k, vv)
		}
	}

	headers.Delete("Content-Length")
	headers.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)

	// create new buffer size 1024
	buf := make([]byte, 32)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				w.WriteChunkedBodyDone()
				break
			}
			fmt.Printf("Error reading httpbin response: %v\n", err)
			return
		}
		n, err = w.WriteChunkedBody(buf[:n])
		if err != nil {
			fmt.Printf("Error writing chunked body: %v\n", err)
			return
		}
		fmt.Printf("Wrote %d bytes\n", n)
	}
}

func writeServerError(w *response.Writer, _ *request.Request) {
	headers := headers.NewHeaders()
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	body := response.BuildResponseBody(response.StatusCodeInternalServerError, "Okay, you know what? This one is on me.")
	w.WriteBody(body)
}

func writeBadRequest(w *response.Writer, _ *request.Request) {
	headers := headers.NewHeaders()
	w.WriteStatusLine(response.StatusCodeBadRequest)
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	body := response.BuildResponseBody(response.StatusCodeBadRequest, "Your request honestly kinda sucked.")
	w.WriteBody(body)
}
