package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/itchyny/gojq"
)

func ParseJsonFile[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var res T
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return res
}

func JsonQuery[T any](input map[string]interface{}, query string) (T, error) {
	var ret T

	parseQuery, err := gojq.Parse(query)
	if err != nil {
		return ret, err
	}

	iter := parseQuery.Run(input) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return ret, err
		}

		if v != nil {
			return v.(T), nil
		}
	}

	return ret, errors.New("Data not found")
}
