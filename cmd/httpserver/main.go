package main

import (
	"log"
	"os"
	"os/signal"
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

	if req.RequestLine.RequestTarget == "/yourproblem" {
		headers := headers.NewHeaders()
		w.WriteStatusLine(response.StatusCodeBadRequest)
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		body := response.BuildResponseBody(response.StatusCodeBadRequest, "Your request honestly kinda sucked.")
		w.WriteBody(body)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		headers := headers.NewHeaders()
		w.WriteStatusLine(response.StatusCodeInternalServerError)
		headers.Set("Content-Type", "text/html")
		w.WriteHeaders(headers)
		body := response.BuildResponseBody(response.StatusCodeInternalServerError, "Okay, you know what? This one is on me.")
		w.WriteBody(body)
		return
	}

	headers := headers.NewHeaders()
	w.WriteStatusLine(response.StatusCodeSuccess)
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	body := response.BuildResponseBody(response.StatusCodeSuccess, "Your request was an absolute banger.")
	w.WriteBody(body)
}
