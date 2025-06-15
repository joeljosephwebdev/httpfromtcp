package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/joeljosephwebdev/httpfromtcp/internal/request"
	"github.com/joeljosephwebdev/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	Port     int
	closed   atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %v", err)
	}
	server := &Server{
		listener: listener,
		Port:     port,
		handler:  handler,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("failed to close server: %v", err)
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("unable to accept connection: %v", err)
			continue
		}
		go s.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	respWriter := response.NewWriter(conn)
	// parse the request from the conn
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			Message:    err.Error(),
			StatusCode: response.StatusCodeBadRequest,
		}
		hErr.Write(respWriter)
		return
	}
	s.handler(respWriter, req)
}

func (he *HandlerError) Write(w *response.Writer) {
	w.WriteStatusLine(he.StatusCode)
	headers := response.GetDefaultHeaders(len(he.Message))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	body := response.BuildResponseBody(he.StatusCode, he.Message)
	w.WriteBody(body)
}
