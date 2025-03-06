package parse

import (
	"reflect"
	//"strings"
	"testing"

	"os"
	"regexp"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

func TestParseToCastFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileContent   string
		expected      *models.CastFile
		expectedError bool
	}{
		{
			name: "Simple Request with Status Assertion",
			fileContent: `
[request]
GET /test HTTP/1.1
Host: example.com

[assert]
status 200
`,
			expected: &models.CastFile{
				CtxMap: map[int]models.HTTPRequestContext{
					0: {
						Request: models.Request{
							Req: nil, // Cannot check fasthttp.Request directly
						},
						Response: models.Response{},
						CmdArgs: models.CommandActions{
							Method:            "GET",
							URL:               "/test",
							Headers:           map[string]string{"Host": "example.com"},
							Body:              []byte{},
							RedirectsToFollow: 10,
							Color:             true,
							VarReplacement:    true,
						},
						Assertions: []models.Assertion{{
							Type:     "status",
							Target:   "",
							Operator: "==",
							Expected: "200",
						}},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Request with Header Assertion",
			fileContent: `
[request]
GET /test HTTP/1.1
Host: example.com

[assert]
header Content-Type == application/json
`,
			expected: &models.CastFile{
				CtxMap: map[int]models.HTTPRequestContext{
					0: {
						Request: models.Request{
							Req: nil, // Cannot check fasthttp.Request directly
						},
						Response: models.Response{},
						CmdArgs: models.CommandActions{
							Method:            "GET",
							URL:               "/test",
							Headers:           map[string]string{"Host": "example.com"},
							Body:              []byte{},
							RedirectsToFollow: 10,
							Color:             true,
							VarReplacement:    true,
						},
						Assertions: []models.Assertion{{
							Type:     "header",
							Target:   "Content-Type",
							Operator: "==",
							Expected: "application/json",
						}},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Chained Request with Variable",
			fileContent: `
[vars]
test_url = /test

[request]
GET {{test_url}} HTTP/1.1
Host: example.com

[assert]
status 200
`,
			expected: &models.CastFile{
				CtxMap: map[int]models.HTTPRequestContext{
					0: {
						Request: models.Request{
							Req: nil, // Cannot check fasthttp.Request directly
						},
						Response: models.Response{},
						CmdArgs: models.CommandActions{
							Method:            "GET",
							URL:               "/test",
							Headers:           map[string]string{"Host": "example.com"},
							Body:              []byte{},
							RedirectsToFollow: 10,
							Color:             true,
							VarReplacement:    true,
						},
						Assertions: []models.Assertion{{
							Type:     "status",
							Target:   "",
							Operator: "==",
							Expected: "200",
						}},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Multiple Assertions",
			fileContent: `
[request]
GET /test HTTP/1.1
Host: example.com

[assert]
status 200
header Content-Type == application/json
body contains test:
`,
			expected: &models.CastFile{
				CtxMap: map[int]models.HTTPRequestContext{
					0: {
						Request: models.Request{
							Req: nil, // Cannot check fasthttp.Request directly
						},
						Response: models.Response{},
						CmdArgs: models.CommandActions{
							Method:            "GET",
							URL:               "/test",
							Headers:           map[string]string{"Host": "example.com"},
							Body:              []byte{},
							RedirectsToFollow: 10,
							Color:             true,
							VarReplacement:    true,
						},
						Assertions: []models.Assertion{
							{
								Type:     "status",
								Target:   "",
								Operator: "==",
								Expected: "200",
							},
							{
								Type:     "header",
								Target:   "Content-Type",
								Operator: "==",
								Expected: "application/json",
							},
							{
								Type:     "body",
								Target:   "contains",
								Operator: "",
								Expected: "test:",
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Two requests",
			fileContent: `
[request]
GET /test HTTP/1.1
Host: example.com

[assert]
status 200

[request]
GET /test2 HTTP/1.1
Host: example.com

[assert]
status 200
`,
			expected: &models.CastFile{
				CtxMap: map[int]models.HTTPRequestContext{
					0: {
						Request: models.Request{
							Req: nil, // Cannot check fasthttp.Request directly
						},
						Response: models.Response{},
						CmdArgs: models.CommandActions{
							Method:            "GET",
							URL:               "/test",
							Headers:           map[string]string{"Host": "example.com"},
							Body:              []byte{},
							RedirectsToFollow: 10,
							Color:             true,
							VarReplacement:    true,
						},
						Assertions: []models.Assertion{{
							Type:     "status",
							Target:   "",
							Operator: "==",
							Expected: "200",
						}},
					},
					1: {
						Request: models.Request{
							Req: nil, // Cannot check fasthttp.Request directly
						},
						Response: models.Response{},
						CmdArgs: models.CommandActions{
							Method:            "GET",
							URL:               "/test2",
							Headers:           map[string]string{"Host": "example.com"},
							Body:              []byte{},
							RedirectsToFollow: 10,
							Color:             true,
							VarReplacement:    true,
						},
						Assertions: []models.Assertion{{
							Type:     "status",
							Target:   "",
							Operator: "==",
							Expected: "200",
						}},
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := &CustomParser{}
			result, err := parser.ParseToCastFile([]byte(tc.fileContent))

			if (err != nil) != tc.expectedError {
				t.Fatalf("expected error: %v, got: %v, error: %v", tc.expectedError, (err != nil), err)
			}

			if err == nil {
				// Clean Request.Req because its impossible to compare it by reflect.DeepEqual
				for key := range result.CtxMap {
					ctx := result.CtxMap[key] // Get a copy of the struct
					ctx.Request.Req = nil     // Modify the copy
					result.CtxMap[key] = ctx  // Assign the copy back to the map
				}

				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("Parsed CastFile does not match expected:\n\nexpected: %+v\n\ngot:      %+v", tc.expected, result)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name           string
		fileContent    string
		expectedCtxMap map[int]models.HTTPRequestContext
		expectedError  bool
	}{
		{
			name: "Parse successful",
			fileContent: `
[request]
GET /test HTTP/1.1
Host: example.com

[assert]
status 200
`,
			expectedCtxMap: map[int]models.HTTPRequestContext{
				0: {
					Request: models.Request{
						Req: nil, // Cannot check fasthttp.Request directly
					},
					Response: models.Response{},
					CmdArgs: models.CommandActions{
						Method:            "GET",
						URL:               "/test",
						Headers:           map[string]string{"Host": "example.com"},
						Body:              []byte{},
						RedirectsToFollow: 10,
						Color:             true,
						VarReplacement:    true,
					},
					Assertions: []models.Assertion{{
						Type:     "status",
						Target:   "",
						Operator: "==",
						Expected: "200",
					}},
				},
			},
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k := koanf.New(".")
			parser := &CustomParser{}
			// Load the configuration.
			err := k.Load(rawbytes.Provider([]byte(tc.fileContent)), parser)
			if (err != nil) != tc.expectedError {
				t.Fatalf("TestParse - expected error: %v, got: %v, error: %v", tc.expectedError, (err != nil), err)
			}

			if err == nil {
				// Access the CastFile.
				castFile := k.Get("castfile").(*models.CastFile)

				if len(castFile.CtxMap) != len(tc.expectedCtxMap) {
					t.Fatalf("TestParse - expected CtxMap length: %d, got: %d", len(tc.expectedCtxMap), len(castFile.CtxMap))
				}

				// Clean Request.Req because its impossible to compare it by reflect.DeepEqual
				for key := range castFile.CtxMap {
					ctx := castFile.CtxMap[key] // Get a copy of the struct
					ctx.Request.Req = nil       // Modify the copy
					castFile.CtxMap[key] = ctx  // Assign the copy back to the map
				}

				if !reflect.DeepEqual(castFile.CtxMap, tc.expectedCtxMap) {
					t.Errorf("TestParse - Parsed CtxMap does not match expected:\n\nexpected: %+v\n\ngot:      %+v", tc.expectedCtxMap, castFile.CtxMap)
				}
			}

		})
	}
}

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		name           string
		fileContent    string
		expectedCtxMap map[int]models.HTTPRequestContext
		expectedError  bool
	}{
		{
			name: "Unmarshal successful",
			fileContent: `
[request]
GET /test HTTP/1.1
Host: example.com

[assert]
status 200
`,
			expectedCtxMap: map[int]models.HTTPRequestContext{
				0: {
					Request: models.Request{
						Req: nil, // Cannot check fasthttp.Request directly
					},
					Response: models.Response{},
					CmdArgs: models.CommandActions{
						Method:            "GET",
						URL:               "/test",
						Headers:           map[string]string{"Host": "example.com"},
						Body:              []byte{},
						RedirectsToFollow: 10,
						Color:             true,
						VarReplacement:    true,
					},
					Assertions: []models.Assertion{{
						Type:     "status",
						Target:   "",
						Operator: "==",
						Expected: "200",
					}},
				},
			},
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := &CustomParser{}
			// Load the configuration.
			data, err := parser.Unmarshal([]byte(tc.fileContent))

			if (err != nil) != tc.expectedError {
				t.Fatalf("TestUnmarshal - expected error: %v, got: %v, error: %v", tc.expectedError, (err != nil), err)
			}

			if err == nil {
				// Access the CastFile.
				castFile := data["castfile"].(*models.CastFile)

				if len(castFile.CtxMap) != len(tc.expectedCtxMap) {
					t.Fatalf("TestUnmarshal - expected CtxMap length: %d, got: %d", len(tc.expectedCtxMap), len(castFile.CtxMap))
				}

				// Clean Request.Req because its impossible to compare it by reflect.DeepEqual
				for key := range castFile.CtxMap {
					ctx := castFile.CtxMap[key] // Get a copy of the struct
					ctx.Request.Req = nil       // Modify the copy
					castFile.CtxMap[key] = ctx  // Assign the copy back to the map
				}

				if !reflect.DeepEqual(castFile.CtxMap, tc.expectedCtxMap) {
					t.Errorf("TestUnmarshal - Parsed CtxMap does not match expected:\n\nexpected: %+v\n\ngot:      %+v", tc.expectedCtxMap, castFile.CtxMap)
				}
			}

		})
	}
}

func TestParseHTTPRequest(t *testing.T) {
	testCases := []struct {
		name           string
		requestString  string
		expectedReq    models.Request
		expectedCmd    models.CommandActions
		expectedError  bool
		expectedHeader map[string]string
	}{
		{
			name: "Simple GET Request",
			requestString: `GET /test HTTP/1.1
Host: example.com`,
			expectedReq: models.Request{
				Req: nil, // Can't check the *fasthttp.Request directly, so skip it
			},
			expectedCmd: models.CommandActions{
				Method:            "GET",
				URL:               "/test",
				Headers:           map[string]string{"Host": "example.com"},
				Body:              []byte{},
				RedirectsToFollow: 10,
				Color:             true,
				VarReplacement:    true,
			},
			expectedError: false,
		},
		{
			name: "POST Request with Body",
			requestString: `POST /test HTTP/1.1
Host: example.com
Content-Type: application/json

{"key": "value"}`,
			expectedReq: models.Request{
				Req: nil, // Can't check the *fasthttp.Request directly, so skip it
			},
			expectedCmd: models.CommandActions{
				Method:  "POST",
				URL:     "/test",
				Headers: map[string]string{"Host": "example.com", "Content-Type": "application/json"},
				Body: []byte(`{"key": "value"}
`), // Note the newline
				RedirectsToFollow: 10,
				Color:             true,
				VarReplacement:    true,
			},
			expectedError: false,
		},
		{
			name:          "Invalid Request",
			requestString: `GET  HTTP/1.1`,
			expectedReq: models.Request{
				Req: nil,
			},
			expectedCmd:   models.CommandActions{},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := &CustomParser{}
			req, cmd, err := parser.parseHTTPRequest(tc.requestString, map[string]string{})

			if (err != nil) != tc.expectedError {
				t.Fatalf("parseHTTPRequest - expected error: %v, got: %v, error: %v", tc.expectedError, (err != nil), err)
			}

			if err == nil {
				// Clear request for comparing
				req.Req = nil

				if !reflect.DeepEqual(req, tc.expectedReq) {
					t.Errorf("parseHTTPRequest - Request does not match expected:\n\nexpected: %+v\n\ngot:      %+v", tc.expectedReq, req)
				}

				if !reflect.DeepEqual(cmd, tc.expectedCmd) {
					t.Errorf("parseHTTPRequest - CommandActions does not match expected:\n\nexpected: %+v\n\ngot:      %+v", tc.expectedCmd, cmd)
				}
			}
		})
	}
}

func TestParseAssertions(t *testing.T) {
	testCases := []struct {
		name            string
		assertionLines  []string
		expected        []models.Assertion
		expectedError   bool
		expectedHeaders map[string]string
	}{
		{
			name: "Status Assertion",
			assertionLines: []string{
				"status 200",
			},
			expected: []models.Assertion{
				{
					Type:     "status",
					Target:   "",
					Operator: "==",
					Expected: "200",
				},
			},
			expectedError: false,
		},
		{
			name: "Header Assertion",
			assertionLines: []string{
				"header Content-Type == application/json",
			},
			expected: []models.Assertion{
				{
					Type:     "header",
					Target:   "Content-Type",
					Operator: "==",
					Expected: "application/json",
				},
			},
			expectedError: false,
		},
		{
			name: "Body Assertion",
			assertionLines: []string{
				"body contains test",
			},
			expected: []models.Assertion{
				{
					Type:     "body",
					Target:   "contains",
					Operator: "", // contains doesn't require operator
					Expected: "test",
				},
			},
			expectedError: false,
		},
		{
			name: "Multiple Assertions",
			assertionLines: []string{
				"status 200",
				"header Content-Type == application/json",
				"body contains test",
			},
			expected: []models.Assertion{
				{
					Type:     "status",
					Target:   "",
					Operator: "==",
					Expected: "200",
				},
				{
					Type:     "header",
					Target:   "Content-Type",
					Operator: "==",
					Expected: "application/json",
				},
				{
					Type:     "body",
					Target:   "contains",
					Operator: "",
					Expected: "test",
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := &CustomParser{}
			assertions, err := parser.parseAssertions(tc.assertionLines)

			if (err != nil) != tc.expectedError {
				t.Fatalf("parseAssertions - expected error: %v, got: %v, error: %v", tc.expectedError, (err != nil), err)
			}

			if err == nil {
				if !reflect.DeepEqual(assertions, tc.expected) {
					t.Errorf("parseAssertions - Assertions does not match expected:\n\nexpected: %+v\n\ngot:      %+v", tc.expected, assertions)
				}
			}
		})
	}
}

func (p *CustomParser) resolveVar(line string, vars map[string]string) string {
	uuidRegex := regexp.MustCompile(`uuid\(\)`)
	envRegex := regexp.MustCompile(`env\.get\("([^"]+)"\)`)
	varRegex := regexp.MustCompile(`{{\s*([a-zA-Z0-9_.]+)\s*}}`)

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

func TestResolveVar(t *testing.T) {
	p := &CustomParser{}
	vars := map[string]string{
		"myvar": "myvalue",
	}
	input := "Hello {{myvar}}"
	expected := "Hello myvalue"
	result := p.resolveVar(input, vars)
	if result != expected {
		t.Errorf("resolveVar - Expected %q, but got %q", expected, result)
	}

	// Skip test case for environment variables in CI environments
	if os.Getenv("CI") != "" {
		t.Skip("Skipping environment variable test in CI environment")
	}

	//envInput := "Hello {{env.get(\"HOME\")}}"
	//envResult := p.resolveVar(envInput, vars)

	//if !strings.Contains(envResult, "/home/") {
	//	t.Errorf("resolveVar - Expected $HOME to be resolved, got: %q", envResult)
	//}
}
