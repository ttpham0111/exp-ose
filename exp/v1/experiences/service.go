package experience

import (
	// "encoding/json"
	"net/http"

	"github.com/ttpham0111/exp-ose/exp/util"
)

type service struct {
}

func (s *service) find(w http.ResponseWriter, r *http.Request) {
	util.JsonResponse(w, "hello", http.StatusOK)
}
