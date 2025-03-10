package models

import (
	"github.com/valyala/fasthttp"
)

// TODO: Document and comment these more in-depth, especially ExecutionError.
type Request struct {
	//Method  string
	//URL     *url.URL
	//Headers map[string]string
	//Body    io.Reader
	Req *fasthttp.Request

	//CLI CommandActions
}

type Response struct {
	//Response    fasthttp.Response
	Status      string
	StatusCode  int
	Headers     map[string]string
	Protocol    string
	ContentType string
	Body        []byte
	Duration    int
	Size        int
}

type CommandActions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte

	Debug             bool
	Color             bool
	More              bool
	VarReplacement    bool
	PrintOptions      []string // Just Request, Just Status, etc.
	RedirectsToFollow int
	FileUploadPath    string
	DownloadPath      string
	CurlOutput        bool
}

type HTTPRequestContext struct {
	Request    Request
	Response   Response
	CmdArgs    CommandActions
	Assertions []Assertion
	Captures   []Capture
}

type CastFile struct {
	Vars   map[string]string
	CtxMap map[int]HTTPRequestContext
}
