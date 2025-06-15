package response

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
)

type responseBodyData struct {
	ResponseBodyTitle   string
	ResponseBodyHeader  string
	ResponseBodyContent string
}

func BuildResponseBody(statusCode StatusCode, content string) []byte {
	tmpl, err := template.ParseFiles("internal/response/templates/response_body.html")
	if err != nil {
		log.Fatal("failed to load template")
	}
	var bodyBuffer bytes.Buffer
	respBody := responseBodyData{
		ResponseBodyTitle:   fmt.Sprintf("%d %s", statusCode, getStatusDescription(statusCode)),
		ResponseBodyHeader:  getStatusDescription(statusCode),
		ResponseBodyContent: content,
	}
	if err := tmpl.Execute(&bodyBuffer, respBody); err != nil {
		log.Fatal("failed to execute template")
	}
	return bodyBuffer.Bytes()
}
