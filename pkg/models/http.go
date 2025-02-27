package models

import (
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// TODO: Document and comment these more in-depth, especially ExecutionError.
type Request struct {
	//Method  string
	//URL     *url.URL
	//Headers map[string]string
	//Body    io.Reader
	Req *fasthttp.Request

	Assertions Assertion
	CLI        CommandActions
}

type Response struct {
	Status      string
	StatusCode  int
	Headers     http.Header
	Protocol    string
	ContentType string
	Body        []byte
	Duration    time.Duration
}

type ExecutionError struct {
	Stage     string // "parsing", "connection", "validation"
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

type CommandActions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte

	Debug             bool
	Highlight         bool
	VarReplacement    bool
	PrintOptions      []string // Just Request, Just Status, etc.
	RedirectsToFollow int
	FileUploadPath    string
}

func (e ExecutionError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Stage, e.Timestamp.Format(time.RFC3339), e.Message)
}
