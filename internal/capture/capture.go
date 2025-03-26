package capture

import (
	"fmt"
	"regexp"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

var GlobalCaptures []models.Capture

// can be grabbed by the koanf parsers
var GlobalVars map[string]string = make(map[string]string) // initializing the map
//var GlobalVars map[string]string

func Capture(ctxHTTP *models.HTTPRequestContext) map[string]string {
	//fmt.Println("reached capture function")
	ctxHTTP.Captures = GlobalCaptures
	//fmt.Println("global capture:")
	//fmt.Println(ctxHTTP.Captures)
	for _, capture := range ctxHTTP.Captures {
		switch {
		case capture.Operation == "header":
			//fmt.Println("reached inside header capture operation")
			err := captureHeaderVal(&ctxHTTP.Response, &capture)
			if err != nil {
				//fmt.Printf("error: %s", err)
				logging.Logger.Error(err)
			}
		case capture.Location == "resp" && capture.Operation == "regex":
			//fmt.Println("reached regex capture.go switch")
			err := captureRegex(&ctxHTTP.Response, &capture)
			if err != nil {
				//fmt.Printf("error: %s", err)
				logging.Logger.Error(err)
			}

		}
	}
	//fmt.Println("reached globalvars being returned capture.go")
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
	//fmt.Println("capture target")
	//fmt.Println(capture.Target)
	re := regexp.MustCompile(capture.Target)
	//fmt.Println(string(resp.Body))
	val := re.Find(resp.Body)
	if val == nil {
		return fmt.Errorf("no value found to capture from body: '%s'", capture.Target)
	}
	GlobalVars[capture.VarName] = string(val)
	//fmt.Println("globalvars from inside regex func")
	//fmt.Println(GlobalVars)
	return nil
}
