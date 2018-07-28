package main

import (
	"os"
)

type config struct {
	port string

	yelpApiKey string
}

func newConfig() *config {
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	return &config{
		port:       port,
		yelpApiKey: os.Getenv("YELP_API_KEY"),
	}
}
