package parse

import (
	"strings"

	"github.com/alp1n3-eth/cast/pkg/models"
)

// Runs fourth
func (p *CustomParser) parseAssertions(assertLines []string) ([]models.Assertion, error) {
	assertions := make([]models.Assertion, 0) // Changed to a slice
	for _, line := range assertLines {
		parts := strings.SplitN(line, " ", 4) // Splitting into 4 parts: Type, Target, Operator, Expected
		if len(parts) < 2 {                   // if parts is less than 2 then skip
			continue // Skip invalid assertion lines - must have all parts.
		}

		if parts[0] == "status" && len(parts) == 2 {
			// Handle status code assertions
			assertion := models.Assertion{
				Type:     "status",
				Target:   "", //status code doesn't have a target
				Operator: "==",
				Expected: parts[1],
			}
			assertions = append(assertions, assertion)
		} else if parts[0] == "body" && len(parts) >= 3 {
			// Handle body contains assertion where the length is greater or equal to three
			assertion := models.Assertion{
				Type:     parts[0],
				Target:   parts[1],                     // contains
				Operator: "",                           // No operator necessary for contains
				Expected: strings.Join(parts[2:], " "), // Joins the other parts in the parts slice into a string
			}
			assertions = append(assertions, assertion)

		} else if len(parts) == 4 {

			assertion := models.Assertion{
				Type:     strings.TrimSpace(parts[0]),
				Target:   strings.TrimSpace(parts[1]),
				Operator: strings.TrimSpace(parts[2]),
				Expected: strings.TrimSpace(parts[3]),
			}
			assertions = append(assertions, assertion)
		} else {
			continue
		}
	}
	return assertions, nil // Returns a slice of assertions
}
