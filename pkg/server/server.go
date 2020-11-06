package server

import (
	"net"
	"sync"
)

const (
	host = "0.0.0.0"
	port = "9999"
)

// HandlerFunc is ...
type HandlerFunc func(conn net.Conn)

// Server is struct
type Server struct {
	addr string
	mu sync.RWMutex
	handlers map[string]HandlerFunc
}

// NewServer is function
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

// Register is method
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

// Start is method
func (s *Server) Start() error {
	
	return nil
}