package models

import (
	"time"
	"net/http"
	"io"
	"net/url"
	"fmt"
)
/*
type HTTPRequestData struct {
	Method string
	URL string
	Headers map[string]string
	Body string
}

type HTTPResponseData struct {
	StatusCode string
	Headers map[string]string
	Body string
}

type HTTPRequest struct {
	Request HTTPRequestData
	Assertions []Assertion
	Response HTTPResponseData
}

type ParsedHTTPFile struct {
	Requests []HTTPRequest
}
 */

// TODO: Document and comment these more in-depth, especially ExecutionError.
type Request struct {
    Method  string
    URL     *url.URL
    Headers http.Header
    Body    io.Reader
}

type Response struct {
    Status 	   string
    StatusCode int
    Headers    http.Header
    Protocol   string
    ContentType string
    Body       []byte
    Duration   time.Duration
    Assertions Assertion
}

type ExecutionError struct {
    Stage     string    // "parsing", "connection", "validation"
    Message   string
    Source    error
    Timestamp time.Time
}

type ExecutionResult struct {
    Request   Request
    Response  Response
    Errors    []ExecutionError
    Timestamp time.Time
}

func (e ExecutionError) Error() string {
    return fmt.Sprintf("[%s] %s: %s", e.Stage, e.Timestamp.Format(time.RFC3339), e.Message)
}
