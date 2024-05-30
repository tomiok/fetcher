package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/tomiok/fetcher/cmd/api"
	"net/http"
	"time"
)

func main() {
	run()
}

func run() {
	// get dependencies.
	deps := api.NewDeps()

	// chi multiplexer.
	r := chi.NewRouter()
	srv := &http.Server{
		Addr: ":9000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// setup http routes.
	api.Routes(r, deps)

	// collect ticker information
	go func() {
		deps.Service.Run(deps.RefreshPeriodInSec)
	}()

	// start web server.
	web := api.Server{
		Server: srv,
	}
	web.Start()
}
