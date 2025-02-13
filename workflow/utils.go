package workflow

import "encoding/json"

func first[T any](value T, _ bool) T {
	return value
}

func pretty[T any](v T) string {
	result, _ := json.Marshal(v)
	return string(result)
}

func prettyIndent[T any](v T) string {
	result, _ := json.MarshalIndent(v, "", "  ")
	return string(result)
}
