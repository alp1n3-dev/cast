package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

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

func Base64(str, cmd string) (string, error) {
	fmt.Println(str + " " + cmd)

	switch cmd {
	case "encode":
		b64str := base64.StdEncoding.EncodeToString([]byte(str))

		return b64str, nil

	case "decode":
		b64byte, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return "", err
		}

		return string(b64byte), nil
	}

	fmt.Println("gonna error out - functions.go")

	return "", fmt.Errorf("FAILURE - Unable to perform Base64 operation")
}
