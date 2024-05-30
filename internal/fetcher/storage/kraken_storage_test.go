//go:build integration

package storage

import (
	"github.com/tomiok/fetcher/internal/fetcher"
	"net/http"
	"testing"
	"time"
)

func TestKrakenStorage_Get(t *testing.T) {
	type fields struct {
		client *fetcher.Kraken
	}
	type args struct {
		pairs []fetcher.Pair
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantQty int
	}{
		{
			name:    "happy path one pair",
			fields:  fields{client: fetcher.NewKraken(&http.Client{}, []string{"BTC/USD", "BTC/EUR"})},
			args:    args{pairs: fetcher.ToPairs([]string{"BTC/USD"})},
			wantQty: 1,
		},
		{
			name:    "happy path two pairs",
			fields:  fields{client: fetcher.NewKraken(&http.Client{}, []string{"BTC/USD", "BTC/EUR"})},
			args:    args{pairs: fetcher.ToPairs([]string{"BTC/USD", "BTC/EUR"})},
			wantQty: 2,
		},
		{
			name:    "wrong pairs",
			fields:  fields{client: fetcher.NewKraken(&http.Client{}, []string{"BTC/USD", "BTC/EUR"})},
			args:    args{pairs: fetcher.ToPairs([]string{"BTC/ARS"})},
			wantQty: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KrakenStorage{
				client: tt.fields.client,
			}
			k.client.Fetch(time.Second * 1) // actually we dont care the time, just once is enough.

			time.Sleep(time.Second * 2) // get some time for the async job to complete.
			got := k.Get(tt.args.pairs)

			if len(got.Ltps) != tt.wantQty {
				t.Fatalf("wrong quantity, got %d, expected %d", len(got.Ltps), tt.wantQty)
			}

			for _, value := range got.Ltps {
				if value.Amount == 0.0 {
					t.FailNow()
				}
			}
		})
	}
}
