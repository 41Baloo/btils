package btils

import (
	"io"

	"github.com/goccy/go-json"
)

// Unmarshal a reader into T and return *T
func Unmarshal[T any](rc io.Reader) (*T, error) {
	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var res T
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Can be slightly faster than 'Unmarshal' since the pointer is passed down
func UnmarshalPointer[T any](in *T, rc io.Reader) (*T, error) {
	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, in)
	if err != nil {
		return nil, err
	}

	return in, nil
}
