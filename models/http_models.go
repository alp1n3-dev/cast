package models

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
