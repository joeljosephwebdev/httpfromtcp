package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type requestState int

const crlf = "\r\n"
const bufferSize = 8

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

func (s requestState) String() string {
	switch s {
	case requestStateInitialized:
		return "requestStateInitialized"
	case requestStateDone:
		return "requestStateDone"
	default:
		return "unknown"
	}
}

type Request struct {
	RequestLine  RequestLine
	RequestState requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var validMethods = map[string]struct{}{
	"GET":     {},
	"POST":    {},
	"PUT":     {},
	"DELETE":  {},
	"HEAD":    {},
	"OPTIONS": {},
	"TRACE":   {},
	"PATCH":   {},
	"CONNECT": {},
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		RequestState: requestStateInitialized,
	}
	for req.RequestState != requestStateDone {
		// if buffer is full, create new buffer twice the size and copy the old data
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.RequestState != requestStateDone {
					return nil, errors.New("request not complete, EOF reached")
				}
				break
			}
			return nil, fmt.Errorf("error reading from request reader: %w", err)
		}
		readToIndex += numBytesRead
		offset, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %w", err)
		}
		// copy the data to a new buffer
		if offset < readToIndex {
			newBuf := make([]byte, readToIndex-offset)
			copy(newBuf, buf[offset:readToIndex])
			buf = newBuf
		} else {
			buf = nil // no data left to copy
		}
		// decrement readtoIndex to account for the offset
		readToIndex -= offset
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + len(crlf), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.RequestState != requestStateInitialized {
		return 0, errors.New("request already parsed or not requestStateInitialized")
	}

	requestLine, offset, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if requestLine == nil && offset == 0 {
		return 0, nil // not enough data to parse request line
	}

	r.RequestLine = *requestLine
	r.RequestState = requestStateDone

	return offset, nil
}

func requestLineFromString(line string) (*RequestLine, error) {
	parts := strings.Fields(line)
	var requestTarget, version, method string

	// check for path
	if len(parts) != 3 {
		return nil, errors.New("malformed request line, must have exactly 3 parts")
	}

	method = parts[0]
	_, ok := validMethods[method]
	if !ok {
		return nil, errors.New("incorrect method used")
	}

	requestTarget = parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", line)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}

	version = versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	requestLine := &RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}

	return requestLine, nil
}
