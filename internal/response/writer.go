package response

import (
	"fmt"

	"github.com/joeljosephwebdev/httpfromtcp/internal/headers"
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.responseState != responseStateInitialized {
		return fmt.Errorf("write status line called out of order")
	}
	defer func() { w.responseState = responseStateHeaders }()
	statusLine := getStatusLine(statusCode)
	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("failed to write status line: %w", err)
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.responseState != responseStateHeaders {
		return fmt.Errorf("write headers called out of order")
	}
	defer func() { w.responseState = responseStateBody }()
	for k, v := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	w.responseState = responseStateBody

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.responseState != responseStateBody {
		return 0, fmt.Errorf("write body called out of order")
	}
	defer func() { w.responseState = responseStateDone }()
	return w.writer.Write(p)
}
