package models

import (
	"bytes"
	"net/http"
)

type ResponseWriterLogWrapper struct {
	W          *http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode *int
}

type RequestResponseLog struct {
	Req        RequestLog  `json:"Req"`
	Resp       ResponseLog `json:"Resp"`
	StatusCode int         `json:"StatusCode"`
}

type RequestLog struct {
	Route  string              `json:"Route"`
	Method string              `json:"Method"`
	Body   string              `json:"Body"`
	Params map[string][]string `json:"Params"`
}

type ResponseLog struct {
	Header map[string][]string `json:"Header"`
	Body   string              `json:"Body"`
}

func (rww ResponseWriterLogWrapper) Write(buf []byte) (int, error) {
	rww.Body.Write(buf)
	return (*rww.W).Write(buf)
}

func (rww ResponseWriterLogWrapper) Header() http.Header {
	return (*rww.W).Header()
}

func (rww ResponseWriterLogWrapper) WriteHeader(statusCode int) {
	(*rww.StatusCode) = statusCode
	(*rww.W).WriteHeader(statusCode)
}

type ResponseWriterWithTimeout struct {
	http.ResponseWriter
	headerWritten bool
}

func (rw *ResponseWriterWithTimeout) WriteHeader(statusCode int) {
	if !rw.headerWritten {
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
	}
}

func (rw *ResponseWriterWithTimeout) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}
