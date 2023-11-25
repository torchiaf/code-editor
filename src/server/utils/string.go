package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func ToString(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(data)
}
