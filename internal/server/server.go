package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/joeljosephwebdev/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	Port     int
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %v", err)
	}
	server := &Server{
		listener: listener,
		Port:     port,
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
	response.WriteStatusLine(conn, response.StatusCodeSucess)
	headers := response.GetDefaultHeaders(0)
	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Printf("failed to write headers: %v", err)
		return
	}
}
