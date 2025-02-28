package models

import (
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
	//Response    fasthttp.Response
	Status      string
	StatusCode  int
	Headers     http.Header
	Protocol    string
	ContentType string
	Body        []byte
	Duration    time.Duration
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
	DownloadPath      string
}

type CastFile struct {
	Variables map[string]string
	Requests  map[*Request]Response
}
