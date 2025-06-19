# httpfromtcp

A minimal HTTP/1.1 server and protocol parser implemented in Go, built from raw TCP sockets. This project demonstrates how HTTP requests can be parsed and served without relying on the standard `net/http` package, providing a learning resource for low-level protocol handling and custom server logic.

> **Note:** This project is guided by the [Boot.dev](https://boot.dev) Deeper Learning curriculum.

---

## Features

- **Custom HTTP/1.1 Request Parsing:**  
  Parses HTTP requests directly from TCP streams, including request line, headers, and body.

- **Flexible Server Architecture:**  
  Modular design with pluggable request handlers and support for custom responses.

- **Chunked Transfer Encoding:**  
  Supports chunked responses and trailers for advanced HTTP scenarios.

- **Static and Dynamic Content Serving:**  
  Serves static files (e.g., MP4 video) and dynamic HTML responses using Go templates.

- **Proxy Functionality:**  
  Includes a simple proxy handler for forwarding requests to [httpbin.org](https://httpbin.org).

- **Comprehensive Testing:**  
  Unit tests for request parsing, header handling, and body extraction.

---

## Getting Started

### Prerequisites

- Go 1.20+ (see [go.mod](go.mod) for details)

### Installation

Clone the repository:

```sh
git clone https://github.com/joeljosephwebdev/httpfromtcp.git
cd httpfromtcp
```

---

## Usage

### Run the TCP Listener

Start the raw TCP listener (for debugging or protocol exploration):

```sh
go run ./cmd/tcplistener | tee /tmp/rawget.http
```

Send a GET request:

```sh
curl http://localhost:42069/coffee
```

Send a POST request:

```sh
curl -X POST -H "Content-Type: application/json" -d '{"flavor":"dark mode"}' http://localhost:42069/coffee
```

### Run the HTTP Server

Start the HTTP server with custom handlers:

```sh
go run ./cmd/httpserver
```

- Visit [http://localhost:42069/](http://localhost:42069/) in your browser.
- Access `/video` to stream the included MP4 file.
- Access `/httpbin/get` to proxy requests to httpbin.org.

### UDP Sender Example

Send UDP packets to a local listener:

```sh
go run ./cmd/udpsender
nc -u -l 42069
```

---

## Project Structure

```
cmd/
  httpserver/      # Main HTTP server entrypoint
  tcplistener/     # Raw TCP listener for HTTP requests
  udpsender/       # UDP sender utility
internal/
  headers/         # HTTP header parsing and utilities
  request/         # HTTP request parsing logic
  response/        # HTTP response construction and templates
  server/          # Server abstraction and handler logic
```

---

## Testing

Run all unit tests:

```sh
go test ./...
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.

---

## Author

Joel Joseph  
[GitHub](https://github.com/joeljosephwebdev)

---

## Acknowledgements

- [Boot.dev](https://boot.dev) for project guidance
- [httpbin.org](https://httpbin.org) for proxy testing
- Go standard library for inspiration
- [RFC 9112: Hypertext Transfer Protocol -- HTTP/1.1](https://datatracker.ietf.org/doc/html/rfc9112)
- [RFC 793: Transmission Control Protocol (TCP)](https://datatracker.ietf.org/doc/html/rfc793)
- [RFC 768: User Datagram Protocol (UDP)](https://datatracker.ietf.org/doc/html/rfc768)
