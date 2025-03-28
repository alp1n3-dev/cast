package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	//"os"
	"regexp"
	"strings"

	"github.com/alp1n3-eth/cast/internal/capture"
	"github.com/alp1n3-eth/cast/internal/utils"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/google/uuid"
)

var vars map[string]string

// Custom parser to handle the specified file format.
type CustomParser struct{}

// Unmarshal implements the koanf.Parser interface.
func (p *CustomParser) Unmarshal(b []byte) (map[string]interface{}, error) {
	return p.Parse(b)
}

// Runs first
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

	// Store variables in local scope instead of global scope.
	vars := make(map[string]string)

	resolveVar := func(line string, vars map[string]string) string {
		if capture.GlobalVars != nil {
			for k, v := range capture.GlobalVars {
				vars[k] = v
			}
		}

		// Perform functions before doing replacements.

		re := regexp.MustCompile(`([a-zA-Z0-9_]+)`) // Regex to capture ANY variable name (alphanumeric and underscore)
		line = re.ReplaceAllStringFunc(line, func(s string) string {
			if val, ok := vars[s]; ok { // Directly use the key "auth_token"
				//fmt.Printf("Replacing '%s' with '%s'", s, val)
				return val
			}
			//fmt.Printf("Variable '%s' not found", s)
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
			case "pre":
				currentSection = "pre"
			case "request":
				// Process previous request if exists
				if len(requestLines) > 0 {
					reqCtx, err := p.parseRequest(requestLines, assertLines, vars, resolveVar) //jsonPathRegex)
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
			case "post":
				currentSection = "post"
			default:
				currentSection = "" // Unknown section
			}
			continue
		}

		switch currentSection {
		case "pre":
			// Run a file before request.
			if !strings.Contains(line, "=") && strings.Contains(line, "run") {
				fmt.Println("[!] running a file of requests before actual provided file")

				parts := strings.Split(line, "(")
				parts2 := strings.Split(parts[1], ")")

				logging.Logger.Infof("Running Request File: %s\n", parts2[0])
				//fmt.Println("target file for pre-run:")
				//fmt.Println(parts2[0])
				//fmt.Println()

				// ./cast file tests/test_files/adv_req_chain.http
				out, err := exec.Command("./cast", "file", parts2[0]).Output()
				if err != nil {
					fmt.Println()
					os.Exit(1)
				}
				//fmt.Println(string(out)) // prints out the output of the file being run.
				logging.Logger.Debugf("File Run Output: %s", out)
				logging.Logger.Info("File Run Complete")
				//fmt.Println("cmd execution finished")
				//os.Exit(0)
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if !strings.HasPrefix(value, `"`) && !strings.HasSuffix(value, `"`) {
					value = runScripts(value)
				}
				value = strings.Trim(value, `"`)
				vars[key] = resolveVar(value, vars)
			}
		case "request":
			if !strings.HasPrefix(line, "#") {
				requestLines = append(requestLines, line)
			}
		case "post":
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
		reqCtx, err := p.parseRequest(requestLines, assertLines, vars, resolveVar) //, jsonPathRegex)
		if err != nil {
			return nil, err
		}
		castFile.CtxMap[requestCounter] = *reqCtx
	}

	return castFile, nil
}

// Marshal marshals the given config map to bytes.
func (p *CustomParser) Marshal(m map[string]interface{}) ([]byte, error) {
	return []byte{}, nil
}

func runScripts(str string) string {
	var value string
	//var err error
	//fmt.Printf("value: %s", value)

	if str == "uuidv7()" {
		uuid, err := uuid.NewV7()
		if err != nil {
			fmt.Println(fmt.Errorf("failed to generate UUID v7"))
			return value
		}
		value = uuid.String()
	}

	if strings.Contains(str, "base64") {
		return base64Ops(str)
	}

	if strings.Contains(str, "url") {
		return urlOps(str)
	}

	return value
}

func urlOps(str string) string {
	var value string
	var err error

	before, after, _ := strings.Cut(str, `"`)
	before, after, _ = strings.Cut(after, `"`)

	if strings.Contains(str, "decode") {
		value, err = url.QueryUnescape(before)
		if err != nil || value == "" {
			//logging.Logger.Fatal(err)
			//fmt.Errorf("%s", err)
			fmt.Println(err)
		}
		return value
	}

	if strings.Contains(str, "encode") {
		value = url.QueryEscape(before)
		if err != nil || value == "" {
			//logging.Logger.Fatal(err)
			//fmt.Errorf("%s", err)
			fmt.Println(err)
		}
		return value
	}

	return value
}

func base64Ops(str string) string {
	var value string
	var err error

	before, after, _ := strings.Cut(str, `"`)
	before, after, _ = strings.Cut(after, `"`)

	if strings.Contains(str, "decode") {
		value, err = utils.Base64(before, "decode")
		if err != nil || value == "" {
			//logging.Logger.Fatal(err)
			//fmt.Errorf("%s", err)
			fmt.Println(err)
		}
		return value
	}

	if strings.Contains(str, "encode") {
		value, err = utils.Base64(before, "encode")
		if err != nil || value == "" {
			//logging.Logger.Fatal(err)
			//fmt.Errorf("%s", err)
			fmt.Println(err)
		}
		return value
	}

	return value
}
