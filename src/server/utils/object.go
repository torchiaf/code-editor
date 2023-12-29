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

func Map[T any](src []T, key func(T) string) map[string]T {
	var result = make(map[string]T)
	for _, v := range src {
		result[key(v)] = v
	}
	return result
}

func Slice[T any](_map map[string]T) []T {
	ret := make([]T, 0, len(_map))
	for _, v := range _map {
		ret = append(ret, v)
	}

	return ret
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func IfNull[T comparable](defValue T, value T) T {
	var nilValue T
	return If(defValue == nilValue, value, defValue)
}
