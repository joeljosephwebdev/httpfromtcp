package headers

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	bytesConsumed := 0

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	//if data begins with a crlf return 2, done, nil
	if idx == 0 {
		return 2, true, nil // indicates end of headers
	}

	// split data by first colon
	colonIndex := bytes.IndexByte(data[:idx], ':')
	if colonIndex == -1 {
		return 0, false, nil // no valid header found
	}

	headerName := string(bytes.ToLower(data[:colonIndex]))
	// if last char of headerName is a space, return 404 error
	if headerName[len(headerName)-1] == ' ' {
		return 0, false, fmt.Errorf("invalid header format: %s", headerName)
	}
	if !validateHeaderName(headerName) {
		return 0, false, fmt.Errorf("invalid header name: %s", headerName)
	}
	headerValue := string(bytes.TrimSpace(data[colonIndex+1 : idx]))

	if len(headerValue) == 0 {
		return 0, false, fmt.Errorf("invalid header value: %s", headerValue)
	}

	if existing, exists := h[headerName]; exists {
		h[headerName] = fmt.Sprintf("%s,%s", existing, headerValue)
	} else {
		h[headerName] = headerValue
	}

	bytesConsumed = idx + len(crlf)
	return bytesConsumed, false, nil
}

func (h Headers) Get(name string) (string, bool) {
	key := strings.ToLower(name)
	value, ok := h[key]
	return value, ok
}

func (h Headers) Add(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Set(key, value string) {
	if !validateHeaderName(strings.ToLower(key)) {
		log.Printf("invalid header name: %s", key)
		return
	}
	h[key] = value
}

func validateHeaderName(name string) bool {
	if len(name) == 0 {
		return false
	}
	// header name can contain any letter or number, or these characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~
	for _, r := range name {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && !bytes.ContainsRune([]byte("!#$%&'*+-.^_`|~"), r) {
			return false
		}
	}
	return true
}
