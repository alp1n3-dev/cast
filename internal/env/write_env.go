package env

import (
	"log"
	"os"
)

func AddEnv(envMap *map[string]string) error {
	// TODO: Write to in-memory .env
	return nil
}

func AddPersistentEnv(kvStr string) error {
	f, err := os.OpenFile("config/persistent.env",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.WriteString("kvStr\n"); err != nil {
		log.Println(err)
	}

	return nil
}

func AddEncryptedKV(envMap *map[string]string) error {
	// TODO: Write to password-protected encrypted persistent .env file.
	return nil
}
