package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

// HandlerFunc ...
type HandlerFunc func(conn net.Conn)

// Server ...
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

// Request is ...
type Request struct {
	Conn        net.Conn
	QueryParams url.Values
}

// NewServer ...
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

// Register ...
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

// Start ...
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			err = cerr
			return
			// if err == nil {
			// log.Print(cerr)
			// }
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

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
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

		requestLine := string(data[:requestLineEnd])
		parts := strings.Split(requestLine, " ")
		if len(parts) != 3 {
			return
		}

		path, version := parts[1], parts[2]

		if version != "HTTP/1.1" {
			return
		}
		decode, err := url.PathUnescape(path)
		if err != nil {
			return
		}
		uri, err := url.ParseRequestURI(decode)
		if err != nil {
			return
		}

		query := uri.Query()
		var req Request
		req.QueryParams = query

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
