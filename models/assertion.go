package models

type Assertion struct {
	Type string // header, header.value, json.name, json.value, body; AKA the type of operation.
	Target string // name of the header, json path,
	Operator string // Such as: ==, >=, <=, !=, ||, &&, <, >, contains
	Expected string // value that is to be operated that is contained in the "Target".
}

/*
TODO: Validate assertions

Assertions won't be able to be validated until AFTER a response body has been received.

Examples:
- Checking for JSON response body if type is set to json.name or json.value.
- Ensuring specific items exist if set as the "Target".
*/

//func (a Assertion) Validate(resp *http.Response) error
