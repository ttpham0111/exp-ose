package services

import (
	"fmt"
)

type ClientError struct {
	Message    string
	StatusCode int
}

func (se ClientError) Error() string {
	return fmt.Sprintf("[%s] %s", se.StatusCode, se.Message)
}
