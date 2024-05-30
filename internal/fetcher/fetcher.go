package fetcher

import (
	"log/slog"
	"os"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

//go:generate mockgen -source=$GOFILE -destination=../mocks/$GOFILE -package=mocks
type Fetcher interface {
	Fetch(d time.Duration)
}

//go:generate mockgen -source=$GOFILE -destination=../mocks/$GOFILE -package=mocks
type Storage interface {
	Get(pairs []Pair) PairResults
}

type Service struct {
	storage Storage

	fetcher Fetcher
}

type Pair struct {
	Key string
}

type LastTradedPrice struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type PairResults struct {
	Ltps []LastTradedPrice `json:"ltp"`
}

func NewService(s Storage, fetcher Fetcher) *Service {
	return &Service{
		storage: s,
		fetcher: fetcher,
	}
}

func (s *Service) GetResults(pairs []Pair) (PairResults, error) {
	return s.storage.Get(pairs), nil
}
func (s *Service) Run(d time.Duration) {
	s.fetcher.Fetch(d)
}

func ToPairs(strPairs []string) []Pair {
	var res []Pair
	for _, pair := range strPairs {
		res = append(res, Pair{Key: pair})
	}

	return res
}
