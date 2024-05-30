package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(r *chi.Mux, deps *Dependencies) {
	r.Use(middleware.Logger, middleware.RequestID, middleware.Recoverer, middleware.Heartbeat("/ping"))

	r.Get("/api/v1/ltp", deps.ltpHandler.GetPrices)
}
