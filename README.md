# Muvor Protocol Secure (MVRPS)

## Description

A secure custom protocol implementation for Muvor Protocol Secure (MVRPS) using TLS.

## Installation
```bash
go get -v github.com/wakuu-enterprises/mvrps-go
```

```bash
go mod tidy
```

## Implementation

## Client

```bash
package main

import (
	"fmt"
	"log"
	"mvrps-protocol/client"
)

func main() {
	req := client.NewSecureRequest("client-key.pem", "client-cert.pem", "ca-cert.pem", "127.0.0.1:8443", "CREATE", "/")
	resp, err := req.Send("Hello, secure server!")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	fmt.Println("Response:", resp)
}
```

## Video

```bash
package main

import (
	"mvrp/mvvp"
)

func main() {
	video.ProcessSegments("/path/to/uploads", "/path/to/structured")
}
```

## Server
```bash
package main

import (
	"log"
	"mvrps-protocol/server"
)

func main() {
	srv := server.NewSecureServer("server-key.pem", "server-cert.pem", "127.0.0.1:8443")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```