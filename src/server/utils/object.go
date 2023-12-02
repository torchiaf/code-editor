package utils

import "github.com/fatih/structs"

func Find[T comparable](items []T, key string, v string) (bool, T) {
	for _, item := range items {
		m := structs.Map(item)
		if m[key] == v {
			return true, item
		}
	}

	return false, *new(T)
}
