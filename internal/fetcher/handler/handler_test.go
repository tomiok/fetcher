package handler

import (
	"fmt"
	"github.com/tomiok/fetcher/internal/fetcher"
	"github.com/tomiok/fetcher/internal/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestHandler_GetPrices(t *testing.T) {
	type fields struct {
		service *fetcher.Service
	}

	ctrl := gomock.NewController(t)
	mockFetcher := mocks.NewMockFetcher(ctrl)
	mockStorage := mocks.NewMockStorage(ctrl)
	service := fetcher.NewService(mockStorage, mockFetcher)

	ltpHandler := &Handler{
		service: service,
	}
	tests := []struct {
		name       string
		fields     fields
		mockFn     func()
		params     []string
		wantStatus int
		wantBody   string
	}{
		{
			name:   "happy path, 2 pairs",
			fields: fields{service: service},
			params: []string{"BTC/USD", "BTC/EUR"},
			mockFn: func() {
				mockStorage.EXPECT().Get(fetcher.ToPairs([]string{"BTC/USD", "BTC/EUR"})).
					Return(fetcher.PairResults{Ltps: []fetcher.LastTradedPrice{
						{
							Pair:   "BTC/USD",
							Amount: "61000",
						},
						{
							Pair:   "BTC/EUR",
							Amount: "62000",
						},
					}})
			},
			wantBody:   `{"ltp":[{"pair":"BTC/USD","amount":"61000"},{"pair":"BTC/EUR","amount":"62000"}]}`,
			wantStatus: http.StatusOK,
		},
		{
			name:   "happy path, 1 pair, other not in our list",
			fields: fields{service: service},
			params: []string{"BTC/USD", "ETH/EUR"},
			mockFn: func() {
				mockStorage.EXPECT().Get(fetcher.ToPairs([]string{"BTC/USD", "ETH/EUR"})).
					Return(fetcher.PairResults{Ltps: []fetcher.LastTradedPrice{
						{
							Pair:   "BTC/USD",
							Amount: "61000",
						},
					}})
			},
			wantBody:   `{"ltp":[{"pair":"BTC/USD","amount":"61000"}]}`,
			wantStatus: http.StatusOK,
		},
		{
			name:   "fails with no one in our list",
			fields: fields{service: service},
			params: []string{"DOGE/EUR", "ETH/EUR"},
			mockFn: func() {
				mockStorage.EXPECT().Get(fetcher.ToPairs([]string{"DOGE/EUR", "ETH/EUR"})).
					Return(fetcher.PairResults{Ltps: []fetcher.LastTradedPrice{}})
			},
			wantBody:   `do not have information for given pairs`,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock interface response.
			tt.mockFn()

			params := strings.Join(tt.params, ",")

			//create request.
			req, err := http.NewRequest(http.MethodGet, "/api/v1/ltp?pairs="+params, nil)
			if err != nil {
				t.FailNow()
			}

			rr := httptest.NewRecorder()
			h := http.HandlerFunc(ltpHandler.GetPrices)

			// execute
			h.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Fatalf("got %d, want %d", rr.Code, tt.wantStatus)
			}

			if !reflect.DeepEqual(strings.TrimSpace(rr.Body.String()), tt.wantBody) {
				fmt.Println(rr.Body.String())
				t.FailNow()
			}
		})
	}
}
