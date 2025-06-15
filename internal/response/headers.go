package response

import (
	"fmt"

	"github.com/joeljosephwebdev/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Add("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Add("Connection", "close")
	h.Add("Content-Type", "text/plain")
	return h
}
