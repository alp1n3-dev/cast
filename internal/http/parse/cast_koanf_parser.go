package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/valyala/fasthttp"
)

// Custom parser to handle the specified file format.
type CustomParser struct{}

// Unmarshal implements the koanf.Parser interface.
func (p *CustomParser) Unmarshal(b []byte) (map[string]interface{}, error) {
	return p.Parse(b)
}

// Parse is internal parse function
func (p *CustomParser) Parse(b []byte) (map[string]interface{}, error) {
	castFile, err := p.ParseToCastFile(b)
	if err != nil {
		return nil, err
	}

	// Convert the CastFile to a map[string]interface{} for koanf compatibility (if needed)
	config := map[string]interface{}{
		"castfile": castFile,
	}

	return config, nil
}

func (p *CustomParser) ParseToCastFile(b []byte) (*models.CastFile, error) {
	castFile := &models.CastFile{
		CtxMap: make(map[int]models.HTTPRequestContext),
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))

	currentSection := ""
	var requestLines []string
	var assertLines []string

	uuidRegex := regexp.MustCompile(`uuid\(\)`)
	envRegex := regexp.MustCompile(`env\.get\("([^"]+)"\)`)
	varRegex := regexp.MustCompile(`{{\s*([a-zA-Z0-9_.]+)\s*}}`)
	jsonPathRegex := regexp.MustCompile(`\$\.([a-zA-Z0-9_]+)`)

	// Store variables in local scope instead of global scope.
	vars := make(map[string]string)

	resolveVar := func(line string, vars map[string]string) string {
		line = uuidRegex.ReplaceAllStringFunc(line, func(s string) string {
			return "generated-uuid" // Or generate a real UUID
		})
		line = envRegex.ReplaceAllStringFunc(line, func(s string) string {
			matches := envRegex.FindStringSubmatch(s)
			if len(matches) > 1 {
				return os.Getenv(matches[1])
			}
			return ""
		})
		line = varRegex.ReplaceAllStringFunc(line, func(s string) string {
			matches := varRegex.FindStringSubmatch(s)
			if len(matches) > 1 {
				varName := matches[1]
				if val, ok := vars[varName]; ok {
					return val
				}
			}
			return s
		})
		return line
	}

	requestCounter := 0
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.Trim(line, "[]")
			switch sectionName {
			case "vars":
				currentSection = "vars"
			case "request":
				// Process previous request if exists
				if len(requestLines) > 0 {
					reqCtx, err := p.parseRequest(requestLines, assertLines, vars, resolveVar, jsonPathRegex)
					if err != nil {
						return nil, err
					}
					castFile.CtxMap[requestCounter] = *reqCtx
					requestCounter++

					// Reset for the new request
					requestLines = []string{}
					assertLines = []string{}
					// Keep vars for chained requests
				}
				currentSection = "request"
			case "assert":
				currentSection = "assert"
			default:
				currentSection = "" // Unknown section
			}
			continue
		}

		switch currentSection {
		case "vars":
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, `"`)
				vars[key] = resolveVar(value, vars)
			}
		case "request":
			if !strings.HasPrefix(line, "#") {
				requestLines = append(requestLines, line)
			}
		case "assert":
			if !strings.HasPrefix(line, "#") {
				assertLines = append(assertLines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Process the last request if exists
	if len(requestLines) > 0 {
		reqCtx, err := p.parseRequest(requestLines, assertLines, vars, resolveVar, jsonPathRegex)
		if err != nil {
			return nil, err
		}
		castFile.CtxMap[requestCounter] = *reqCtx
	}

	return castFile, nil
}

func (p *CustomParser) parseRequest(
	requestLines []string,
	assertLines []string,
	vars map[string]string,
	resolveVar func(line string, vars map[string]string) string,
	jsonPathRegex *regexp.Regexp,
) (*models.HTTPRequestContext, error) {
	requestStr := strings.Join(requestLines, "\n")

	// Preserve the newlines within the JSON body
	var finalRequestLines []string
	inJSON := false
	for _, requestLine := range strings.Split(requestStr, "\n") {
		if strings.Contains(requestLine, "{") {
			inJSON = true
		}
		if inJSON {
			finalRequestLines = append(finalRequestLines, requestLine)
		} else {
			finalRequestLines = append(finalRequestLines, strings.TrimSpace(requestLine))
		}
		if strings.Contains(requestLine, "}") {
			inJSON = false
		}
	}

	requestStr = strings.Join(finalRequestLines, "\n")

	requestStr = resolveVar(requestStr, vars)
	requestStr = jsonPathRegex.ReplaceAllStringFunc(requestStr, func(s string) string {
		return s // Placeholder for jsonpath replacement logic later.
	})

	req, cmdArgs, err := p.parseHTTPRequest(requestStr, vars)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTTP request: %w", err)
	}

	assertions, err := p.parseAssertions(assertLines)
	if err != nil {
		return nil, fmt.Errorf("error parsing assertions: %w", err)
	}

	reqCtx := &models.HTTPRequestContext{
		Request:    req,
		Response:   models.Response{}, // Initialize an empty response
		CmdArgs:    cmdArgs,
		Assertions: assertions,
	}
	fmt.Println(reqCtx)

	return reqCtx, nil
}

func (p *CustomParser) parseHTTPRequest(requestStr string, vars map[string]string) (models.Request, models.CommandActions, error) {
	// Split the request string into lines.
	lines := strings.Split(requestStr, "\n")

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
	fmt.Println(parsedURL)
	fmt.Println(urlStr)
	fmt.Println(fullURL)
	if err != nil && urlStr != "" {
		return models.Request{}, models.CommandActions{}, fmt.Errorf("invalid URL: %w", err)
	}

	// Parse headers and body.
	headers := make(map[string]string)
	var bodyBuilder strings.Builder
	isBody := false

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if line == "" {
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
			}
		}
	}

	body := bodyBuilder.String()

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
		RedirectsToFollow: 10,
		Color:             true,
		VarReplacement:    true,
	}

	request := models.Request{
		Req: req,
	}

	fmt.Println(request)
	fmt.Println(cmdActions)

	return request, cmdActions, nil
}

func (p *CustomParser) parseAssertions(assertLines []string) ([]models.Assertion, error) {
	assertions := make([]models.Assertion, 0) // Changed to a slice
	for _, line := range assertLines {
		parts := strings.SplitN(line, " ", 4) // Splitting into 4 parts: Type, Target, Operator, Expected
		if len(parts) < 2 {                   // if parts is less than 2 then skip
			continue // Skip invalid assertion lines - must have all parts.
		}

		if parts[0] == "status" && len(parts) == 2 {
			// Handle status code assertions
			assertion := models.Assertion{
				Type:     "status",
				Target:   "", //status code doesn't have a target
				Operator: "==",
				Expected: parts[1],
			}
			assertions = append(assertions, assertion)
		} else if parts[0] == "body" && len(parts) >= 3 {
			// Handle body contains assertion where the length is greater or equal to three
			assertion := models.Assertion{
				Type:     parts[0],
				Target:   parts[1],                     // contains
				Operator: "",                           // No operator necessary for contains
				Expected: strings.Join(parts[2:], " "), // Joins the other parts in the parts slice into a string
			}
			assertions = append(assertions, assertion)

		} else if len(parts) == 4 {

			assertion := models.Assertion{
				Type:     strings.TrimSpace(parts[0]),
				Target:   strings.TrimSpace(parts[1]),
				Operator: strings.TrimSpace(parts[2]),
				Expected: strings.TrimSpace(parts[3]),
			}
			assertions = append(assertions, assertion)
		} else {
			continue
		}
	}
	fmt.Println(assertions)
	return assertions, nil // Returns a slice of assertions
}

// Marshal marshals the given config map to bytes.
func (p *CustomParser) Marshal(m map[string]interface{}) ([]byte, error) {
	return []byte{}, nil
}
