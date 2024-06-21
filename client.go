package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mvrp-protocol/mvrp"
)

type SecureRequest struct {
	addr    string
	method  string
	url     string
	config  *tls.Config
}

func NewSecureRequest(keyFile, certFile, caFile, addr, method, url string) *SecureRequest {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return &SecureRequest{
		addr:   addr,
		method: method,
		url:    url,
		config: config,
	}
}

func (r *SecureRequest) Send(body string) (string, error) {
	conn, err := tls.Dial("tcp", r.addr, r.config)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %v", r.addr, err)
	}
	defer conn.Close()

	request := fmt.Sprintf("%s %s MVRP/1.0\r\nContent-Length: %d\r\n\r\n%s", r.method, r.url, len(body), body)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	response := string(buffer[:n])
	return response, nil
}
