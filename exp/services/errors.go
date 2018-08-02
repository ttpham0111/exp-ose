package services

import (
	"fmt"
)

type ClientError struct {
	JsonResponse string
	Message      string
	StatusCode   int
}

func (err ClientError) Error() string {
	msg := err.JsonResponse

	if msg == "" {
		msg = fmt.Sprintf("[%v] %v", err.StatusCode, err.Message)
	}

	return msg
}
