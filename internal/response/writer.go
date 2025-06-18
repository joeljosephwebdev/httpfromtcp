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

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.responseState != responseStateBody {
		return 0, fmt.Errorf("write body called out of order")
	}
	defer func() { w.responseState = responseStateDone }()
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.responseState != responseStateBody {
		return 0, fmt.Errorf("write body called out of order")
	}
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.responseState != responseStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.responseState)
	}
	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	w.responseState = responseStateTrailers
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.responseState != responseStateTrailers {
		return fmt.Errorf("cannot write trailers in state %d", w.responseState)
	}
	defer func() { w.responseState = responseStateBody }()
	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
