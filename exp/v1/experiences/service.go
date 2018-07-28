package experience

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ttpham0111/exp-ose/exp"
)

type service struct {
}

func (s *service) find(w http.ResponseWriter, r *http.Request) {
	exp.JsonResponse(w, []], http.StatusOK)
}

func (s *service) get(w http.ResponseWriter, r *http.Request) {
	experienceId := chi.URLParam(r, "experienceId")
}
