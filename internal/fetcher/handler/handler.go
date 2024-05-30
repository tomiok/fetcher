package handler

import (
	"encoding/json"
	"github.com/tomiok/fetcher/internal/fetcher"
	"net/http"
	"strings"
)

type Handler struct {
	service *fetcher.Service
}

func NewHandler(service *fetcher.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("pairs")
	param = strings.ToUpper(param)

	pairs := strings.Split(param, ",")

	res, err := h.service.GetResults(fetcher.ToPairs(pairs))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(res.Ltps) == 0 {
		http.Error(w, "do not have information for given pairs", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
