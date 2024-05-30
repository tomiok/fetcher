package api

import (
	"github.com/tomiok/fetcher/internal/fetcher"
	"github.com/tomiok/fetcher/internal/fetcher/handler"
	"github.com/tomiok/fetcher/internal/fetcher/storage"
	"net/http"
	"time"
)

var currencies = []string{
	"BTC/USD", "BTC/EUR", "BTC/CHF",
}

const refreshPeriodIn = time.Second * 30

type Dependencies struct {
	Service *fetcher.Service

	ltpHandler         *handler.Handler
	RefreshPeriodInSec time.Duration
}

func NewDeps() *Dependencies {
	httpClient := http.Client{}
	kraken := fetcher.NewKraken(&httpClient, currencies)
	stg := storage.NewKrakenStorage(kraken)
	service := fetcher.NewService(stg, kraken)

	ltpHandler := handler.NewHandler(service)
	return &Dependencies{
		Service:            service,
		ltpHandler:         ltpHandler,
		RefreshPeriodInSec: refreshPeriodIn,
	}
}
