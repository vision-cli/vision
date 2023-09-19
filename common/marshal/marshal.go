package marshal

import (
	"bytes"
	"encoding/json"
)

func Unmarshal[T any](data string) (T, error) {
	var resp T
	dec := json.NewDecoder(bytes.NewReader([]byte(data)))
	dec.DisallowUnknownFields() // Force errors
	err := dec.Decode(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func Marshal[T any](data T) (string, error) {
	respStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(respStr), nil
}
