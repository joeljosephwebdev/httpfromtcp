package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

const crlf = "\r\n"

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

	headerName := string(data[:colonIndex])
	// if last char of headerName is a space, return 404 error
	if headerName[len(headerName)-1] == ' ' {
		return 0, false, errors.New("invalid header format")
	}
	headerValue := string(bytes.TrimSpace(data[colonIndex+1 : idx]))

	if len(headerValue) == 0 {
		return 0, false, errors.New("invalid header value")
	}
	h[headerName] = headerValue
	bytesConsumed = idx + len(crlf)

	return bytesConsumed, false, nil

}
