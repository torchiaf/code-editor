package utils

import (
	"log"
	"os"
)

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return data
}
