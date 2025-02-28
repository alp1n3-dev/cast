package env

import (
	"log"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

// For reading and storing environment variables.
// Main pkg: https://github.com/knadh/koanf
// We are using the koanf file provider and dotenv parser.
func ReadEnv() (*map[string]string, error) {
	var envMap map[string]string

	return &envMap, nil
}

func ReadKVFile(filepath string) (*map[string]string, error) {
	var envMap map[string]string

	// Placeholder for testing
	//filepath = "../../config/persistent.env"

	if err := k.Load(file.Provider(filepath), dotenv.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	//fmt.Println(k.String("test")) retrieved test=value from the file.
	envMap = k.StringMap("") // Reads entire thing
	//fmt.Println(envMap)

	return &envMap, nil
}

func ReadEncryptedKV(password string) (*map[string]string, error) {
	var envMap map[string]string

	return &envMap, nil
}
