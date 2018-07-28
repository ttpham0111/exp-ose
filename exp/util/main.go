package util

import (
	"encoding/json"
	"log"
	"net/http"
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

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	switch d := data.(type) {
	case string:
		if _, err := w.Write([]byte(d)); err != nil {
			panic(err)
		}
	default:
		json.NewEncoder(w).Encode(data)
	}
}
