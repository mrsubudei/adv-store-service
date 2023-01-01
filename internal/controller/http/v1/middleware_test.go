package v1_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrsubudei/adv-store-service/internal/config"
	v1 "github.com/mrsubudei/adv-store-service/internal/controller/http/v1"
	"github.com/mrsubudei/adv-store-service/internal/entity"
	mock "github.com/mrsubudei/adv-store-service/internal/service/mock"
	"github.com/mrsubudei/adv-store-service/pkg/logger"
)

func setup() *v1.Handler {
	cfg, err := config.LoadConfig("../../../../config.json")
	if err != nil {
		log.Fatal(err)
	}
	l := logger.New()

	mockService := mock.NewMockService()
	handler := v1.NewHandler(mockService, cfg, l)
	handler.NewRouteGroups()

	return handler
}

func getMockHandler() http.HandlerFunc {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		limit := 0
		offset := 0
		sortBy := ""
		orderBy := ""
		if val, ok := r.Context().Value(entity.KeyLimit).(int); ok && val != 0 {
			limit = val
		}
		if val, ok := r.Context().Value(entity.KeyOffset).(int); ok && val != 0 {
			offset = val
		}
		if val, ok := r.Context().Value(entity.KeySortBy).(string); ok && val != "" {
			if val == "price" {
				sortBy = val
			}
		}
		if val, ok := r.Context().Value(entity.KeyOrderBy).(string); ok && val != "" {
			orderBy = val
		}
		if limit != 0 && offset != 0 && sortBy != "" && orderBy != "" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	})
	return mockHandler
}

func TestParseQuery(t *testing.T) {
	handler := setup()

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantResult string
	}{
		{
			name:       "OK",
			url:        "/v1/adverts?limit=10&offset=20&sort_by=price&order_by=asc",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Error wrong query: fields",
			url:        "/v1/adverts?fields=abc",
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"queries have wrong value","detail":"'fields=' query value should be 'true'"}`,
		},
		{
			name:       "Error wrong query: sort_by",
			url:        "/v1/adverts?sort_by=abc",
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"queries have wrong value","detail":"'sort_by=' query value should be either 'created_at' or 'price'"}`,
		},
		{
			name:       "Error wrong query: order_by",
			url:        "/v1/adverts?order_by=abc",
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"queries have wrong value","detail":"'order_by=' query value should be either 'asc' or 'desc'"}`,
		},
		{
			name:       "Error wrong query: offset",
			url:        "/v1/adverts?offset=abc",
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"queries have wrong value","detail":"'offset=' query value should be positive number"}`,
		},
		{
			name:       "Error wrong query: limit",
			url:        "/v1/adverts?limit=abc",
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"queries have wrong value","detail":"'limit=' query value should be positive number"}`,
		},
		{
			name:       "Error wrong query: limit",
			url:        "/v1/adverts?limit=abc",
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"queries have wrong value","detail":"'limit=' query value should be positive number"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockHandler := getMockHandler()
			handlerToTest := handler.ParseQuery(mockHandler)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)

			handlerToTest.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			} else if rec.Body.String() != tt.wantResult {
				t.Fatalf("want: %v, got: %v", tt.wantResult, rec.Body.String())
			}
		})
	}
}
