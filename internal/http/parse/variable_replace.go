package parse

import (
	//"bufio"
	//"fmt"
	//"os"
	"regexp"
	"strings"

	//"github.com/alp1n3-eth/cast/pkg/logging"
	//"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/valyala/fasthttp"
)

func SwapReqVals(req *fasthttp.Request, replacementVariables *map[string]string) {
	// Regex pattern to match placeholders like {{ FUZZ }} or {{FUZZ}}
	//placeholderPattern := regexp.MustCompile(`{{\s*(\w+)\s*}}`)
	//fmt.Println("reached swapreqvals")

	// Replace in URI
	uri := string(req.URI().FullURI())
	//ri = placeholderPattern.ReplaceAllStringFunc(uri, func(match string) string {
	//key := strings.TrimSpace(match[2 : len(match)-2]) // Extract key inside {{ }}
	//if value, exists := (*replacementVariables)[key]; exists {
	//return value
	//}
	//return match // Leave unchanged if no replacement found
	//})
	req.SetRequestURI(replaceAllVariables(uri, replacementVariables))

	// Replace in Headers
	req.Header.VisitAll(func(key, value []byte) {
		headerVal := string(value)
		//headerVal = placeholderPattern.ReplaceAllStringFunc(headerVal, func(match string) string {
		//key := strings.TrimSpace(match[2 : len(match)-2])
		//if value, exists := (*replacementVariables)[key]; exists {
		//return value
		//}
		//return match
		//})
		headerVal = replaceAllVariables(headerVal, replacementVariables)
		logging.Logger.Debugf("Swapping Req Vals Key: %s, Value: %s", key, headerVal)
		req.Header.SetBytesKV(key, []byte(headerVal))
	})

	// Replace in Body
	body := string(req.Body())
	//body = placeholderPattern.ReplaceAllStringFunc(body, func(match string) string {
	//key := strings.TrimSpace(match[2 : len(match)-2])
	//if value, exists := (*replacementVariables)[key]; exists {
	//return value
	//}
	//return match
	//})

	req.SetBody([]byte(replaceAllVariables(body, replacementVariables)))
}

func replaceAllVariables(target string, replacementVariables *map[string]string) string {
	placeholderPattern := regexp.MustCompile(`{{\s*(\w+)\s*}}`)

	result := placeholderPattern.ReplaceAllStringFunc(target, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		if value, exists := (*replacementVariables)[key]; exists {

			return value
		}
		return match
	})

	return result
}
