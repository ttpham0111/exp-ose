package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/ttpham0111/exp-ose/exp/v1"
)

func newRouter(conf *config) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Blueprints
	v1Blueprint := v1.Blueprint{
		yelpApiKey: conf.yelpApiKey,
	}
	router.Route("/api/v1", v1Blueprint.Register)

	return router
}

func main() {
	conf := newConfig()
	router := newRouter(conf)

	log.Fatal(http.ListenAndServe(":"+conf.port, router))
}
