package assert

import (
	"strings"

	//"github.com/alp1n3-eth/cast/pkg/logging"
	//"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
)

func CheckIfJSONVarExists(jsonBody *map[string]interface{}, targetStr *string) (bool, error) {
	var exists bool

	return exists, nil
}

func RetrieveJSONValue(resp *fasthttp.Response, targetStr *string) (string, error) {
	var retrievedStr string

	return retrievedStr, nil
}

// Run this first.
func checkIfBodyIsJSON(resp *fasthttp.Response) (bool, *map[string]interface{}) {
	var bodyValidity bool
	var jsonBody map[string]interface{}

	// Check content-type header
	if !(strings.Contains(string(resp.Header.Peek("Content-Type")), "application/json")) {
		return false, nil
	}

	// Check the body
	//if err := sonic.Unmarshal(resp.Body(), jsonBody); err != nil {
	//logging.Logger.Fatalf("Unable to unmarshal JSON: %s, %s", resp.Body(), err)
	//return false, nil
	//}

	return bodyValidity, &jsonBody
}
