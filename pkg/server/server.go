package server

import (
	"strings"
	"bytes"
	"io"
	"log"
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
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err != nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn)
	}
}

// handle is method
func (s *Server) handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != io.EOF {
			log.Printf("%s", buf[:n])
			return
		}
		if err != nil {
			return
		}
		data := buf[:n]
		requestLineDelim := []byte{'\r', '\n'}
		requestLineEnd := bytes.Index(data, requestLineDelim)
		if requestLineEnd == -1 {
			return
		}
		requesLine := string(data[:requestLineEnd])
		parts := strings.Split(requesLine, " ")
		if len(parts) != 3 {
			return
		}
		path, version := parts[1], parts[2]
		if version != "HTTP/1.1" {
			return
		}
		handler := func(conn net.Conn) {
			err := conn.Close()
			if err != nil {
				log.Print(err)
			}
		}
		s.mu.RLock()
		for i := 0; i < len(s.handlers); i++ {
			if handl, ok := s.handlers[path]; ok {
				handler = handl
				break
			}
		}
		s.mu.RUnlock()
		handler(conn)
	}
}