package assert

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

func ValidateAssertions(resp *models.Response, assertions []models.Assertion, result *models.ResultOut) {
	var err error
	logging.Logger.Debug(assertions)
	for _, assertion := range assertions {
		switch {
		case assertion.Type == "status":
			//expectedStatus := strings.TrimPrefix(assertion.Type, "status ")
			err = validateStatusCode(resp, &assertion.Expected)
			if err != nil {
				logging.Logger.Error(err)
				result.FailureTotal += 1
			}

		case assertion.Type == "header.value":
			// Header value check
			err = validateHeaderValue(resp, &assertion)
			if err != nil {
				logging.Logger.Error(err)
				result.FailureTotal += 1
			}

		/*
			// Header presence/absence checks
			err = validateHeader(resp, &assertion)
			if err != nil {
				logging.Logger.Error(err)
			}
		*/
		case assertion.Type == "header":
			//fmt.Println(assertion)
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
				result.FailureTotal += 1
			}

		case assertion.Type == "regex":
		// regex

		case assertion.Type == "size":
			err = validateSize(resp, &assertion)
			if err != nil {
				logging.Logger.Error(err)
				result.FailureTotal += 1
			}
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

	logging.Logger.Info("status assertion successful")

	return nil
}

func validateHeaderValue(resp *models.Response, assertion *models.Assertion) error {

	if assertion.Operator == "!=" {
		if assertion.Expected == resp.Headers[assertion.Target] {
			return fmt.Errorf("NEGATIVE ASSERT FAILURE - Not Expected: %s, Actual: %s", assertion.Target, assertion.Expected)
		}
		return nil
	}

	if assertion.Expected == resp.Headers[assertion.Target] {
		logging.Logger.Infof("ASSERT SUCCESS - Expected: %s, Actual: %s", assertion.Target, assertion.Expected)
		return nil
	}

	return fmt.Errorf("ASSERT FAILURE - Expected: %s, Target: %s", assertion.Expected, assertion.Target)
}

func validateHeader(resp *models.Response, assertion *models.Assertion) error {
	//headerContents := resp.Header.PeekAll(assertion.Target)
	if assertion.Operator == "!=" {
		for _, f := range resp.Headers {
			if strings.Contains(f, assertion.Target) {

				return fmt.Errorf("header assertion failed. Expect '%s' to NOT be present", assertion.Expected)
			}
		}

		logging.Logger.Infof("ASSERT SUCCESS - Expected: %s, is NOT present", assertion.Expected)
		return nil
	}

	//fmt.Println(resp.Headers)
	//fmt.Println(assertion.Target)
	//fmt.Println(assertion.Expected)

	if value, ok := resp.Headers[assertion.Expected]; ok {
		logging.Logger.Infof("ASSERT SUCCESS - Expected: %s, Actual: %s", assertion.Expected, value)
		return nil
	}

	return fmt.Errorf("ASSERT FAILURE - Expected: %s", assertion.Expected)
}

func validateBody(resp *models.Response, expectedStr string) error {
	if bytes.Contains(resp.Body, []byte(expectedStr)) {
		logging.Logger.Infof("ASSERT SUCCESS - Expected: %s, Location: Body", expectedStr)
		return nil
	}

	return fmt.Errorf("ASSERT FAILURE - Expected: %s, Location: Body", expectedStr)
}

func validateRegex(resp *models.Response, expectedStr string) error {

	return nil
}

func validateSize(resp *models.Response, assertion *models.Assertion) error {
	expected, err := strconv.Atoi(assertion.Expected)
	if err != nil {
		return fmt.Errorf("error attempting to convert expected to int. Expected: %s", assertion.Expected)
	}

	if assertion.Operator == ">" {
		if resp.Size > expected {
			logging.Logger.Infof("ASSERT SUCCESS - Expected Size: %d", expected)
			return nil
		}
	}

	if assertion.Operator == ">=" {
		if resp.Size >= expected {
			logging.Logger.Infof("ASSERT SUCCESS - Expected Size: %d", expected)
			return nil
		}
	}

	if assertion.Operator == "<" {
		if resp.Size < expected {
			logging.Logger.Infof("ASSERT SUCCESS - Expected Size: %d", expected)
			return nil
		}
	}

	if assertion.Operator == "<=" {
		if resp.Size <= expected {
			logging.Logger.Infof("ASSERT SUCCESS - Expected Size: %d", expected)
			return nil
		}
	}

	if assertion.Operator == "==" {
		if resp.Size == expected {
			logging.Logger.Infof("ASSERT SUCCESS - Expected Size: %d", expected)
			return nil
		}
	}

	if assertion.Operator == "!=" {
		if resp.Size != expected {
			logging.Logger.Infof("ASSERT SUCCESS - Expected Size: %d", expected)
			return nil
		}
	}

	return fmt.Errorf("ASSERT FAILURE - Expected: %s, Operator: %s", assertion.Expected, assertion.Operator)
}

func validateJSON(resp *models.Response, expectedStr string) error {
	//

	return nil
}
