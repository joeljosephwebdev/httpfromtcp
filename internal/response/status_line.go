package response

import (
	"fmt"
)

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) string {
	reasonPhrase := getStatusDescription(statusCode)
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)
}

func getStatusDescription(statusCode StatusCode) string {
	switch statusCode {
	case StatusCodeSuccess:
		return "OK"
	case StatusCodeBadRequest:
		return "Bad Request"
	case StatusCodeInternalServerError:
		return "Internal Server Error"
	default:
		return "Internal Server Error"
	}
}
