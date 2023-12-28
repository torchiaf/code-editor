package utils

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes/scheme"

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

func ParseK8sResource[T any](path string) T {
	decode := scheme.Codecs.UniversalDeserializer().Decode

	stream, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	obj, _, err := decode(stream, nil, nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return obj.(T)
}
