package main

import (
	"time"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"bytes"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Print(err)
		return err
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

		err = handle(conn)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
		return nil
	}
	if err != nil {
		return err
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

	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET"{
		return
	}

	if version != "HTTP/1.1"{
		return
	}

	if path == "/" {
		body, err := ioutil.ReadFile("static/index.html")
		if err != nil {
			return fmt.Errorf("can't read index.html: %w", err)
		}
		
		marker := "{{year}}"
		year := time.Now().Year()
		body = bytes.ReplaceAll(body,[]byte(marker), []byte(strconv.Itoa(year)))

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
			"Content-Type: text/html\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			string(body),
		))
		if err != nil {
			return err
		}
	}
	return nil
}
