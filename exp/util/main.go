package util

import (
	"log"
	"os"

	"gopkg.in/go-playground/validator.v8"
)

func Getenv(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}

func EnsureEnv(name string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	log.Fatal("Missing environmental variable: " + name)
	return ""
}

func FirstError(errors validator.ValidationErrors) *validator.FieldError {
	for _, err := range errors {
		return err
	}
	return nil
}
