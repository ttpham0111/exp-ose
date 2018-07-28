package experience

import (
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	experiences := service{}

	router.Get("/experiences", experiences.find)
	router.Get("/experiences/{experienceId}", experiences.get)

	return router
}
