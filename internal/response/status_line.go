package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusCodeSucess              StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) string {
	reasonPhrase := ""
	switch statusCode {
	case StatusCodeSucess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		reasonPhrase = "Internal Server Error"
	}
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write([]byte(getStatusLine(statusCode)))
	if err != nil {
		return fmt.Errorf("failed to write status line: %w", err)
	}
	return nil
}
