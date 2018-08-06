package main

import (
	"github.com/ttpham0111/exp-ose/exp/util"
)

type config struct {
	port string

	dbUrl  string
	dbName string

	yelpApiKey         string
	eventbriteApiToken string
}

func newConfig() *config {
	port := util.Getenv("PORT", "3000")

	return &config{
		port: port,

		dbUrl:  util.EnsureEnv("DB_URL"),
		dbName: util.EnsureEnv("DB_NAME"),

		yelpApiKey:         util.EnsureEnv("YELP_API_KEY"),
		eventbriteApiToken: util.EnsureEnv("EVENTBRITE_API_TOKEN"),
	}
}
