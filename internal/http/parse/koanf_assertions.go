package parse

import (
	"strings"

	"github.com/alp1n3-eth/cast/internal/http/capture"
	"github.com/alp1n3-eth/cast/pkg/models"
)

// Runs fourth
func (p *CustomParser) parseAssertions(assertLines []string) ([]models.Assertion, error) {
	assertions := make([]models.Assertion, 0) // Changed to a slice
	//captures := make([]models.Capture, 0)
	//captures := &capture.GlobalCaptures

	for _, line := range assertLines {
		parts := strings.SplitN(line, " ", 4) // Splitting into 4 parts: Type, Target, Operator, Expected
		if len(parts) < 2 {                   // if parts is less than 2 then skip
			continue // Skip invalid assertion lines - must have all parts.
		}

		//if len(parts[1]) == 1 {
		//logging.Logger.Debug("CAPTURE DETECTED CAPTURE DETECTED")
		//}
		//if strings.HasSuffix(, suffix string)
		for i := range parts {
			parts[i] = strings.TrimSuffix(parts[i], `"`)
			parts[i] = strings.TrimPrefix(parts[i], `"`)
		}

		if parts[1] == "=" && len(parts) == 4 {
			//fmt.Printf("REACHED PARSING CAPTURE")
			//parts[3] = strings.TrimSuffix(parts[3], `"`)

			//parts[3] = strings.TrimPrefix(parts[3], `"`)

			capture1 := models.Capture{
				Location:  "resp",
				Target:    parts[3],
				VarName:   parts[0],
				Operation: parts[2],
			}

			capture.GlobalCaptures = append(capture.GlobalCaptures, capture1)
			//fmt.Println(capture.GlobalCaptures)
			//captures = append(captures, capture1)

		}

		if parts[0] == "status" {
			// Handle status code assertions
			assertion := models.Assertion{
				Type:     "status",
				Target:   "", //status code doesn't have a target
				Operator: parts[1],
				Expected: parts[2],
			}
			assertions = append(assertions, assertion)
		} else if parts[0] == "header" {
			/*
				if len(parts) == 2 {
					// Handle status code assertions
					assertion := models.Assertion{
						Type:     "header",
						Target:   "", //status code doesn't have a target
						Operator: "==",
						Expected: parts[1],
					}
					assertions = append(assertions, assertion)
				}
			*/
			if parts[0] == "header" && len(parts) == 4 {
				assertion := models.Assertion{
					Type:     "header.value", // header
					Target:   parts[1],       // target header
					Operator: parts[2],       // ==, !=
					Expected: parts[3],       // wanted value
				}
				//fmt.Println(assertion)
				assertions = append(assertions, assertion)
			}

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
