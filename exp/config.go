package main

type config struct {
	port string

	dbUrl  string
	dbName string

	yelpApiKey string
}

func newConfig() *config {
	port := Getenv("PORT", "3000")

	return &config{
		port: port,

		dbUrl:  EnsureEnv("DB_URL"),
		dbName: EnsureEnv("DB_NAME"),

		yelpApiKey: EnsureEnv("YELP_API_KEY"),
	}
}
