package assert

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/valyala/fasthttp"
)

func ValidateAssertions(resp *models.Response, assertions []models.Assertion) {
	var err error
	logging.Logger.Debug(assertions)
	for _, assertion := range assertions {
		switch {
		case assertion.Type == "status":
			//expectedStatus := strings.TrimPrefix(assertion.Type, "status ")
			err = validateStatusCode(resp, &assertion.Expected)
			if err != nil {
				logging.Logger.Error(err)
			}

		case assertion.Type == "header":
			// Header presence/absence checks
			err = validateHeader(resp, &assertion)
			if err != nil {
				logging.Logger.Error(err)
			}

		case assertion.Type == "body":
			// body
			err = validateBody(resp, assertion.Expected)
			if err != nil {
				logging.Logger.Error(err)
			}

		case assertion.Type == "regex":
		// regex

		case assertion.Type == "size":
		// size

		case assertion.Type == "json":
			// json body

		}
	}
}

func validateStatusCode(resp *models.Response, expectedStr *string) error {
	expectedInt, err := strconv.Atoi(*expectedStr)
	if err != nil {
		err = fmt.Errorf("unable to retrieve expected value for assertion")
		return err
	}

	if resp.StatusCode != expectedInt {
		return fmt.Errorf("status assertion failed. expected: %d, actual: %d", expectedInt, resp.StatusCode)
	}

	return nil
}

func validateHeader(resp *models.Response, assertion *models.Assertion) error {
	//headerContents := resp.Header.PeekAll(assertion.Target)
	if assertion.Operator == "NOT" {
		if strings.Contains(resp.Headers, assertion.Target) {
			return fmt.Errorf("header assertion failed. Expect '%s' to NOT be present", assertion.Expected)
		}
		return nil
	}

	if strings.Contains(resp.Headers, assertion.Target) {
		if strings.Contains(resp.Headers, assertion.Expected) {
			return nil
		}
	}

	return fmt.Errorf("header assertion failed. Expect '%s' to be present", assertion.Expected)
	/*
		if headerContents == nil {
			return fmt.Errorf("header validation target does not exist")
		}

		searchBytes := []byte(assertion.Expected)
		for _, row := range headerContents {
			if bytes.Contains(row, searchBytes) {
				return nil
			}
		}
		return nil
	*/
}

func validateBody(resp *models.Response, expectedStr string) error {
	if bytes.Contains(resp.Body, []byte(expectedStr)) {
		return nil
	}

	return fmt.Errorf("body assertion failed")
}

func validateRegex(resp *fasthttp.Response, expectedStr string) error {

	return nil
}

func validateSize(resp *fasthttp.Response, expectedInt int) error {
	//

	return nil
}

func validateJSON(resp *fasthttp.Response, expectedStr string) error {
	//

	return nil
}
