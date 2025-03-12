package parse

import (
	//"fmt"
	"strings"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/valyala/fasthttp"
)

func SwapReqVals(req *fasthttp.Request, replacementVariables *map[string]string) {
	// Regex pattern to match placeholders like {{ FUZZ }} or {{FUZZ}}
	//placeholderPattern := regexp.MustCompile(`{{\s*(\w+)\s*}}`)
	//fmt.Println("reached swapreqvals")

	// Replace in URI
	uri := string(req.URI().FullURI())
	req.SetRequestURI(replaceAllVariables(uri, replacementVariables))

	// Replace in Headers
	req.Header.VisitAll(func(key, value []byte) {
		headerVal := string(value)
		headerVal = replaceAllVariables(headerVal, replacementVariables)
		logging.Logger.Debugf("Swapping Req Vals Key: %s, Value: %s", key, headerVal)
		req.Header.SetBytesKV(key, []byte(headerVal))
	})

	// Replace in Body
	body := string(req.Body())
	req.SetBody([]byte(replaceAllVariables(body, replacementVariables)))
}

func replaceAllVariables(target string, replacementVariables *map[string]string) string {
	/*
		placeholderPattern := regexp.MustCompile(`{{\s*(\w+)\s*}}`)


		result := placeholderPattern.ReplaceAllStringFunc(target, func(match string) string {
			key := strings.TrimSpace(match[2 : len(match)-2])
			if value, exists := (*replacementVariables)[key]; exists {

				return value
			}
			return match
		})

		return result
	*/

	for key, value := range *replacementVariables {
		if strings.Contains(target, key) {
			//fmt.Println("replacing with strings.contains, key=", key, "value=", value)
			logging.Logger.Debugf("Replacing [%s] using key [%s] & value [%s]", target, key, value)

			target = strings.ReplaceAll(target, key, value)
		}
	}
	return target
}
