package models

type HTTPRequestData struct {
	Method string
	URL string
	Headers map[string]string
	Body string
}

type HTTPRequest struct {
	Request HTTPRequestData
	Assertions []Assertion
}

type Assertion struct {
	Type string
	Value string
}

type ParsedHTTPFile struct {
	Requests []HTTPRequest
}
