package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/valyala/fasthttp"
)

// Runs second
func (p *CustomParser) parseRequest(
	requestLines []string,
	assertLines []string,
	vars map[string]string,
	resolveVar func(line string, vars map[string]string) string,
) (*models.HTTPRequestContext, error) {
	requestStr := strings.Join(requestLines, "\n")

	// Preserve the newlines within the JSON body
	var finalRequestLines []string
	//inJSON := false
	for _, requestLine := range strings.Split(requestStr, "\n") {

		finalRequestLines = append(finalRequestLines, strings.TrimSpace(requestLine))
	}

	requestStr = strings.Join(finalRequestLines, "\n")

	requestStr = resolveVar(requestStr, vars)

	req, cmdArgs, err := p.parseHTTPRequest(requestStr, vars)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP request: %w", err)
	}

	//logging.Logger.Debug(req)

	assertions, err := p.parseAssertions(assertLines)
	if err != nil {
		return nil, fmt.Errorf("error parsing assertions: %w", err)
	}

	logging.Logger.Debug(assertions)

	reqCtx := &models.HTTPRequestContext{
		Request:    req,
		Response:   models.Response{}, // Initialize an empty response
		CmdArgs:    cmdArgs,
		Assertions: assertions,
	}

	//logging.Logger.Debug(reqCtx)

	return reqCtx, nil
}

// Runs third
func (p *CustomParser) parseHTTPRequest(requestStr string, vars map[string]string) (models.Request, models.CommandActions, error) {
	// Split the request string into lines.
	lines := strings.Split(requestStr, "\n")
	//fmt.Println(lines)
	//fmt.Println("lines above")

	// Extract method and URL from the first line.
	methodURL := strings.SplitN(lines[0], " ", 2)
	if len(methodURL) < 2 {
		return models.Request{}, models.CommandActions{}, fmt.Errorf("invalid request line: %s", lines[0])
	}
	method := methodURL[0]
	fullURL := ""
	if len(methodURL) > 1 {
		fullURL = methodURL[1]
	}

	logging.Init(true)

	// Split the full URL by space to get only the URL path
	urlParts := strings.SplitN(fullURL, " ", 2)
	if len(urlParts) == 0 {
		// Modified: Empty URL is not an error if there was no URL provided originally
		if fullURL != "" {
			return models.Request{}, models.CommandActions{}, fmt.Errorf("invalid URL: %s", fullURL) // check for empty
		}
		urlParts = []string{""} // Handle case with no URL
	}
	urlStr := urlParts[0] // Takes only the path

	parsedURL, err := url.Parse(urlStr)

	logging.Logger.Debug(fullURL)

	if err != nil && urlStr != "" {
		return models.Request{}, models.CommandActions{}, fmt.Errorf("invalid URL: %w", err)
	}

	// Parse headers and body.
	headers := make(map[string]string)
	var bodyBuilder strings.Builder
	isBody := false

	//logging.Logger.Error(lines)

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if line == "" && !isBody {
			// Empty line indicates the start of the body.
			isBody = true
			continue
		}

		if isBody {
			bodyBuilder.WriteString(line)
			bodyBuilder.WriteString("\n")
		} else {
			headerParts := strings.SplitN(line, ":", 2)
			if len(headerParts) == 2 {
				key := strings.TrimSpace(headerParts[0])
				value := strings.TrimSpace(headerParts[1])
				headers[key] = value
			} else {
				isBody = true
				bodyBuilder.WriteString(line)
				bodyBuilder.WriteString("\n")
			}
		}
	}

	body := bodyBuilder.String()

	//logging.Logger.Debug(body)

	// Construct fasthttp request
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(parsedURL.String())
	req.Header.SetMethod(method)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	hostHeader := req.Header.Peek("Host")

	if len(hostHeader) < 4 {
		return models.Request{}, models.CommandActions{}, fmt.Errorf("invalid host in host header: %w", err)
	}

	req.SetBody([]byte(body))

	// Construct command actions
	cmdActions := models.CommandActions{
		Method:  method,
		URL:     parsedURL.String(),
		Headers: headers,
		Body:    []byte(body),
		// Set other fields to default values, can be overridden later.
		//RedirectsToFollow: 10,
		//Color:             true,
		VarReplacement: true,
	}

	request := models.Request{
		Req: req,
	}

	logging.Logger.Debug(request)
	//logging.Logger.Debug(cmdActions)

	return request, cmdActions, nil
}
