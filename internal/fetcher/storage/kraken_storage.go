package storage

import (
	"github.com/tomiok/fetcher/internal/fetcher"
)

type KrakenStorage struct {
	client *fetcher.Kraken
}

func NewKrakenStorage(k *fetcher.Kraken) *KrakenStorage {
	return &KrakenStorage{client: k}
}

func (k *KrakenStorage) Get(pairs []fetcher.Pair) fetcher.PairResults {
	return k.client.GetValues(pairs)
}
