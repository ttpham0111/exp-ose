package main

import (
	"github.com/ttpham0111/exp-ose/exp/util"
)

type config struct {
	port string

	yelpApiKey string
}

func newConfig() *config {
	port := util.Getenv("PORT", "3000")

	return &config{
		port:       port,
		yelpApiKey: util.EnsureEnv("YELP_API_KEY"),
	}
}
