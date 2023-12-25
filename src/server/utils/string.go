package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(n int, r ...string) string {

	var alphabet []rune

	if len(r) > 0 {
		alphabet = []rune(r[0])
	} else {
		alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	}

	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}

func ToString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(data)
}
