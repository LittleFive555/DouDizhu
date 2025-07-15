package network

import "fmt"

const (
	DeserializeError = 1000
	SerializeError   = 1001
)

type InnerError struct {
	Code    int
	Message string
}

func (e *InnerError) Error() string {
	return fmt.Sprintf("InnerError: %v", e.Message)
}
