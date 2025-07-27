package utils

import "encoding/json"

func ToJsonByte[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

func ToJsonString[T any](v T) (string, error) {
	b, err := ToJsonByte(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
