package services

import (
	"fmt"
)

type ClientError struct {
	Message    string
	StatusCode int
}

func (err ClientError) Error() string {
	return fmt.Sprintf("[%s] %s", err.StatusCode, err.Message)
}
