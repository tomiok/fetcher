package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const baseURL = "https://api.kraken.com/0/public/Ticker?pair=%s"

type Kraken struct {
	Pairs      []Pair
	httpClient *http.Client

	cache *sync.Map
}

type Response struct {
	Error  []interface{}         `json:"error"`
	Result map[string]TickerInfo `json:"result"`
}

type TickerInfo struct {
	C [2]string `json:"c"`
}

func NewKraken(httpClient *http.Client, currencies []string) *Kraken {
	return &Kraken{
		Pairs:      ToPairs(currencies),
		httpClient: httpClient,
		cache:      &sync.Map{},
	}
}

func (k *Kraken) Fetch(d time.Duration) {
	for _, c := range k.Pairs {
		go func(key string) {
			for {
				res, err := k.httpClient.Get(fmt.Sprintf(baseURL, key))
				if err != nil {
					logger.Warn("cannot fetch value", "pair", key, "err", err)
					continue
				}
				var (
					body      = res.Body
					krakenRes = Response{}
				)

				if err := json.NewDecoder(body).Decode(&krakenRes); err != nil {
					logger.Warn("cannot fetch value", "pair", key, "err", err)
					continue
				}

				if len(krakenRes.Error) > 0 {
					logger.Warn("cannot fetch value asd", "pair", key, "err", err)
					continue
				}

				logger.Info("info fetched", "res", krakenRes)

				k.cache.Store(key, resolveAmount(key, krakenRes))
				<-time.After(d)
			}

		}(c.Key)
	}
}

func (k *Kraken) GetValues(pairs []Pair) PairResults {
	var res PairResults
	for _, pair := range pairs {
		value, ok := k.cache.Load(pair.Key)
		if !ok {
			logger.Error("cannot load value", "pair", pair)
			continue
		}

		amount := value.(string)
		fAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			logger.Error("cannot parse amount", "err", err, "pair", pair)
			continue
		}

		res.Ltps = append(res.Ltps, LastTradedPrice{
			Pair:   pair.Key,
			Amount: fAmount,
		})
	}

	return res
}

func resolveAmount(key string, v any) string {
	res, ok := v.(Response)
	if !ok {
		return ""
	}
	tickerInfo := res.Result[key]
	return tickerInfo.C[0]
}
