package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/mvrp-protocol/mvrp"
)

type SecureServer struct {
	addr    string
	config  *tls.Config
}

func NewSecureServer(keyFile, certFile, addr string) *SecureServer {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return &SecureServer{
		addr:   addr,
		config: config,
	}
}

func (s *SecureServer) ListenAndServe() error {
	listener, err := tls.Listen("tcp", s.addr, s.config)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.addr, err)
	}
	defer listener.Close()

	log.Printf("MVRPS server listening on %s", s.addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *SecureServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("failed to read from connection: %v", err)
		return
	}

	request := string(buffer[:n])
	lines := strings.Split(request, "\r\n")
	if len(lines) < 1 {
		log.Printf("malformed request: %s", request)
		return
	}

	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) < 2 {
		log.Printf("malformed request line: %s", requestLine)
		return
	}

	method := parts[0]
	url := parts[1]
	headers := make(map[string]string)
	body := ""

	for i, line := range lines[1:] {
		if line == "" {
			body = strings.Join(lines[i+2:], "\r\n")
			break
		}
		headerParts := strings.SplitN(line, ": ", 2)
		if len(headerParts) == 2 {
			headers[headerParts[0]] = headerParts[1]
		}
	}

	log.Printf("Received %s request for %s with body: %s", method, url, body)
	s.handleRequest(conn, method, url, headers, body)
}

func (s *SecureServer) handleRequest(conn net.Conn, method, url string, headers map[string]string, body string) {
	responseHeaders := map[string]string{
		"Content-Type": "text/plain",
	}

	var responseBody string
	var statusLine string

	switch method {
	case "OPTIONS":
		statusLine = "MVRP/1.0 204 No Content"
		responseHeaders["Allow"] = "OPTIONS, CREATE, READ, EMIT, BURN"
	case "CREATE":
		statusLine = "MVRP/1.0 201 Created"
		responseBody = "Resource created\n"
	case "READ":
		statusLine = "MVRP/1.0 200 OK"
		responseBody = "Resource read\n"
	case "EMIT":
		statusLine = "MVRP/1.0 200 OK"
		responseBody = "Event emitted\n"
	case "BURN":
		statusLine = "MVRP/1.0 200 OK"
		responseBody = "Resource burned\n"
	default:
		statusLine = "MVRP/1.0 405 Method Not Allowed"
		responseBody = "Method not allowed\n"
	}

	response := fmt.Sprintf("%s\r\n", statusLine)
	for k, v := range responseHeaders {
		response += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	response += "\r\n" + responseBody

	conn.Write([]byte(response))
}
