package capture

import (
	"fmt"
	"regexp"

	"github.com/alp1n3-eth/cast/pkg/models"
)

var GlobalCaptures []models.Capture

// can be grabbed by the koanf parsers
var GlobalVars map[string]string

func Capture(ctxHTTP *models.HTTPRequestContext) map[string]string {
	fmt.Println("reached capture function")
	ctxHTTP.Captures = GlobalCaptures
	for _, capture := range ctxHTTP.Captures {
		switch {
		case capture.Location == "header":
			captureHeaderVal(&ctxHTTP.Response, &capture)
		case capture.Location == "resp" && capture.Operation == "regex":
			fmt.Println("reached regex capture.go switch")
			captureRegex(&ctxHTTP.Response, &capture)

		}
	}
	fmt.Println("reached globalvars being returned capture.go")
	return GlobalVars
}

func captureHeaderVal(resp *models.Response, capture *models.Capture) error {
	if val, ok := resp.Headers[capture.Target]; ok {
		GlobalVars[capture.VarName] = val
		return nil
	}
	return fmt.Errorf("no value found to capture for header: '%s'", capture.Target)
}

func captureRegex(resp *models.Response, capture *models.Capture) error {
	re := regexp.MustCompile(capture.Target)
	val := re.Find(resp.Body)
	if val == nil {
		return fmt.Errorf("no value found to capture from body: '%s'", capture.Target)
	}
	GlobalVars[capture.VarName] = string(val)
	return nil
}
