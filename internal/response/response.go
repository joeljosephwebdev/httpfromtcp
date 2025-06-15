package response

import (
	"io"
)

type Writer struct {
	writer        io.Writer
	responseState ResponseState
}

type ResponseState int

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer:        w,
		responseState: responseStateInitialized,
	}
}

const (
	responseStateInitialized ResponseState = iota
	responseStateHeaders
	responseStateBody
	responseStateDone
)
