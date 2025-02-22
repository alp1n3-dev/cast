package parse

import (
	//"bufio"
	//"fmt"
	//"os"
	"strings"
	"regexp"

	//"github.com/alp1n3-eth/cast/pkg/logging"
	//"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/valyala/fasthttp"
)

func SwapReqVals(req *fasthttp.Request, replacementVariables *map[string]string) {
	// Regex pattern to match placeholders like {{ FUZZ }} or {{FUZZ}}
		placeholderPattern := regexp.MustCompile(`{{\s*(\w+)\s*}}`)

		// Replace in URI
		uri := string(req.URI().FullURI())
		uri = placeholderPattern.ReplaceAllStringFunc(uri, func(match string) string {
			key := strings.TrimSpace(match[2 : len(match)-2]) // Extract key inside {{ }}
			if value, exists := (*replacementVariables)[key]; exists {
				return value
			}
			return match // Leave unchanged if no replacement found
		})
		req.SetRequestURI(uri)

		// Replace in Headers
		req.Header.VisitAll(func(key, value []byte) {
			headerVal := string(value)
			headerVal = placeholderPattern.ReplaceAllStringFunc(headerVal, func(match string) string {
				key := strings.TrimSpace(match[2 : len(match)-2])
				if value, exists := (*replacementVariables)[key]; exists {
					return value
				}
				return match
			})
			req.Header.SetBytesKV(key, []byte(headerVal))
		})

		// Replace in Body
		body := string(req.Body())
		body = placeholderPattern.ReplaceAllStringFunc(body, func(match string) string {
			key := strings.TrimSpace(match[2 : len(match)-2])
			if value, exists := (*replacementVariables)[key]; exists {
				return value
			}
			return match
		})
		req.SetBody([]byte(body))
}
