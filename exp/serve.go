package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/ttpham0111/exp-ose/exp/util"
	"github.com/ttpham0111/exp-ose/exp/v1"
)

const (
	version = "0.0.1"
)

func newRouter(conf *config) *chi.Mux {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Blueprints
	v1Blueprint := v1.Blueprint{
		YelpApiKey: conf.yelpApiKey,
	}

	router.Route("/v1", v1Blueprint.Register)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		util.JsonResponse(
			w,
			map[string]string{
				"status":  "OK",
				"version": version,
			},
			http.StatusOK,
		)
	})

	return router
}

func printRoutes(router *chi.Mux) {
	walkFunc := func(
		method string,
		route string,
		_ http.Handler,
		_ ...func(http.Handler) http.Handler,
	) error {
		log.Printf("  %s %s\n", method, strings.Replace(route, "/*", "", -1))
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conf := newConfig()
	router := newRouter(conf)

	log.Println("Serving on localhost:" + conf.port)
	printRoutes(router)

	log.Fatal(http.ListenAndServe(":"+conf.port, router))
}
