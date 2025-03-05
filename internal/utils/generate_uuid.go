package utils

import "github.com/google/uuid"

func GenerateUUID(v string) (string, error) {
	var id string
	//var err error

	if v == "7" {
		id7, err := uuid.NewV7()
		if err != nil {
			return id, err
		}
		id = id7.String()
		return id, nil
	}

	id = uuid.New().String()

	return id, nil
}
