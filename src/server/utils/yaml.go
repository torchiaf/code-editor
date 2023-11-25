package utils

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseFile[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var res T
	err = yaml.Unmarshal([]byte(data), &res)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return res
}
