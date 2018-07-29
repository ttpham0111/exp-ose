package main

import (
	"log"
	"os"
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
