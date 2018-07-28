package experience

import (
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	experiences := service{}

	router.Get("/", experiences.find)

	return router
}
