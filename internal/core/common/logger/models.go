package logger

import (
	"bytes"
	"net/http"
)

type responseWriterWrapper struct {
	W          *http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode *int
}

type requestResponseLog struct {
	Req        requestLog  `json:"Req"`
	Resp       responseLog `json:"Resp"`
	StatusCode int         `json:"StatusCode"`
}

type requestLog struct {
	Route  string              `json:"Route"`
	Method string              `json:"Method"`
	Body   string              `json:"Body"`
	Params map[string][]string `json:"Params"`
}

type responseLog struct {
	Header map[string][]string `json:"Header"`
	Body   string              `json:"Body"`
}
