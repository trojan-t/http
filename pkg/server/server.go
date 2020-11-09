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

// HandlerFunc is func . . .
type HandlerFunc func(req *Request)

// Server is struct of server . . .
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

// Request is struct
type Request struct {
	Conn        net.Conn
	QueryParams url.Values
	PathParams  map[string]string
	Headers     map[string]string
	Body        []byte
}

// NewServer creates new server struct
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

// Register is method . . .
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

// Start is method for start Server
func (s *Server) Start() (err error) {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			err = cerr
			return
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handle(conn)
	}
}

// handle is method . . .
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, (1024 * 8))
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("%s", buf[:n])
		}
		if err != nil {
			log.Println(err)
			return
		}
		var rqst Request
		data := buf[:n]
		rLD := []byte{'\r', '\n'}
		rLE := bytes.Index(data, rLD)
		if rLE == -1 {
			log.Printf("Bad Request")
			return
		}
		headLD := []byte{'\r', '\n', '\r', '\n'}
		headLE := bytes.Index(data, headLD)
		if rLE == -1 {
			return
		}
		headersLine := string(data[rLE:headLE])
		headers := strings.Split(headersLine, "\r\n")[1:]
		mp := make(map[string]string)
		for _, v := range headers {
			headerLine := strings.Split(v, ": ")
			mp[headerLine[0]] = headerLine[1]
		}
		rqst.Headers = mp
		bdy := string(data[headLE:])
		bdy = strings.Trim(bdy, "\r\n")
		rqst.Body = []byte(bdy)
		rqstLine := string(data[:rLE])
		parts := strings.Split(rqstLine, " ")
		if len(parts) != 3 {
			log.Print("Bad Request")
			return
		}
		path, version := parts[1], parts[2]
		if version != "HTTP/1.1" {
			log.Print("HTTP version must be 1.1.")
			return
		}
		decode, err := url.PathUnescape(path)
		if err != nil {
			log.Println(err)
			return
		}
		uri, err := url.ParseRequestURI(decode)
		if err != nil {
			log.Println(err)
			return
		}
		rqst.Conn = conn
		rqst.QueryParams = uri.Query()
		var handler = func(req *Request) { conn.Close() }
		s.mu.RLock()
		pParam, hr := s.checkPath(uri.Path)
		if hr != nil {
			handler = hr
			rqst.PathParams = pParam
		}
		s.mu.RUnlock()
		handler(&rqst)
	}
}

func (s *Server) checkPath(path string) (map[string]string, HandlerFunc) {
	strRoutes := make([]string, len(s.handlers))
	i := 0
	for k := range s.handlers {
		strRoutes[i] = k
		i++
	}
	mp := make(map[string]string)
	for i := 0; i < len(strRoutes); i++ {
		flag := false
		route := strRoutes[i]
		partsRoute := strings.Split(route, "/")
		pRotes := strings.Split(path, "/")
		for j, v := range partsRoute {
			if v != "" {
				f := v[0:1]
				l := v[len(v)-1:]
				if f == "{" && l == "}" {
					mp[v[1:len(v)-1]] = pRotes[j]
					flag = true
				} else if pRotes[j] != v {
					strs := strings.Split(v, "{")
					if len(strs) > 0 {
						key := strs[1][:len(strs[1])-1]
						mp[key] = pRotes[j][len(strs[0]):]
						flag = true
					} else {
						flag = false
						break
					}
				}
				flag = true
			}
		}
		if flag {
			if hr, found := s.handlers[route]; found {
				return mp, hr
			}
			break
		}
	}
	return nil, nil
}
