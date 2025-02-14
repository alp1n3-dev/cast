package assert

import (
	"github.com/alp1n3-eth/cast/models"
)

func AssertController (reqresp models.HTTPRequest) error {

	// TODO: Determine how many asserts there are for the request.
	// Run through each. If one fails, continue. If there is no response, hard failure.

	// TODO: Check if there is a valid response that was saved alongside the originating request.
	validResponse(reqresp)
	// TODO:

	return nil
}

func validResponse (reqresp models.HTTPRequest) bool {
	if len(reqresp.Response.Body) > 5 {
		return true
	}

	return false
}
