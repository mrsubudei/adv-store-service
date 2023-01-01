package v1_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

var (
	advert1 = entity.Advert{
		Name:         "first item",
		Description:  "asd",
		Price:        40,
		MainPhotoUrl: "http://files.com/12",
		PhotosUrls: []string{
			"http://files.com/12",
			"http://files.com/13",
		},
	}
	advert2 = entity.Advert{
		Name:         "second item",
		Description:  "dgdrg",
		Price:        50,
		MainPhotoUrl: "http://files.com/14",
		PhotosUrls: []string{
			"http://files.com/14",
			"http://files.com/16",
		},
	}
)

func TestCreateAdvert(t *testing.T) {
	handler := setup()

	tests := []struct {
		name       string
		reqData    string
		wantStatus int
		wantResult string
	}{
		{
			name:       "OK",
			reqData:    `{"name":"car","description":"asd","price":40,"photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantStatus: http.StatusCreated,
			wantResult: `{"data":[{"id":1}]}`,
		},
		{
			name: "Error name already exist",
			reqData: `{"name":"car","description":"asd","price":40,
			"photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantStatus: http.StatusConflict,
			wantResult: `{"error":"item with name 'car' already exists"}`,
		},
		{
			name:       "Error wrong data format",
			reqData:    `{"name":5}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"json format is not correct"}`,
		},
		{
			name: "Error too many url links",
			reqData: `{"name":"first item","description":"asd","price":40,
			"photo_urls":["http://files.com/12","http://files.com/13", "http://files.com/14",
			"http://files.com/15"]}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"Request Entity Too Large","detail":"'photo_urls:' field's quantity exceeded"}`,
		},
		{
			name: "Error empty field",
			reqData: `{"description":"asd","price":40,
			"photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"request has empty fields","detail":"'name:' field is required"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/adverts", bytes.NewReader([]byte(tt.reqData)))
			handler.Mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			} else if rec.Body.String() != tt.wantResult {
				t.Fatalf("want: %v, got: %v", tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestGetAllAdverts(t *testing.T) {
	handler := setup()
	ctx := context.Background()

	t.Run("OK status no content", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/adverts", nil)
		handler.Mux.ServeHTTP(rec, req)

		wantResult := `{}`

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		} else if rec.Body.String() != wantResult {
			t.Fatalf("want: %v, got: %v", wantResult, rec.Body.String())
		}
	})

	handler.Service.Create(ctx, advert1)
	handler.Service.Create(ctx, advert2)

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v1/adverts", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusAccepted, rec.Code)
		}
	})
}

func TestGetAdvert(t *testing.T) {
	handler := setup()
	handler.Service.Create(context.Background(), advert1)
	handler.Service.Create(context.Background(), advert2)

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantResult string
	}{
		{
			name:       "OK",
			wantStatus: http.StatusOK,
			url:        "/v1/adverts/2",
			wantResult: `{"data":[{"name":"second item","price":50,"main_photo_url":"http://files.com/14"}]}`,
		},
		{
			name:       "OK with additional fields",
			wantStatus: http.StatusOK,
			url:        "/v1/adverts/2?fields=true",
			wantResult: `{"data":[{"name":"second item","description":"dgdrg","price":50,"main_photo_url":"http://files.com/14","photo_urls":["http://files.com/14","http://files.com/16"]}]}`,
		},
		{
			name:       "Error does not exist",
			url:        "/v1/adverts/5",
			wantStatus: http.StatusNotFound,
			wantResult: `{"error":"no content found with id: 5"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			handler.Mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			} else if rec.Body.String() != tt.wantResult {
				t.Fatalf("want: %v, got: %v", tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestUpdateAdvert(t *testing.T) {
	handler := setup()
	handler.Service.Create(context.Background(), advert1)

	tests := []struct {
		name       string
		reqData    string
		url        string
		wantStatus int
		wantResult string
	}{
		{
			name:       "OK",
			wantStatus: http.StatusOK,
			url:        "/v1/adverts/1",
			reqData:    `{"name":"new name","photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantResult: `{}`,
		},
		{
			name:       "Error item does not exist",
			wantStatus: http.StatusNotFound,
			url:        "/v1/adverts/5",
			reqData:    `{"name":"new name","photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantResult: `{"error":"no content found with id: 5"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewReader([]byte(tt.reqData)))
			handler.Mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			} else if rec.Body.String() != tt.wantResult {
				t.Fatalf("want: %v, got: %v", tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestDeleteAdvert(t *testing.T) {
	handler := setup()
	handler.Service.Create(context.Background(), advert1)

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantResult string
	}{
		{
			name:       "OK",
			wantStatus: http.StatusOK,
			url:        "/v1/adverts/1",
			wantResult: `{}`,
		},
		{
			name:       "Error item does not exist",
			wantStatus: http.StatusNotFound,
			url:        "/v1/adverts/5",
			wantResult: `{"error":"no content found with id: 5"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			handler.Mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			} else if rec.Body.String() != tt.wantResult {
				t.Fatalf("want: %v, got: %v", tt.wantResult, rec.Body.String())
			}
		})
	}
}

func TestCommonGroup(t *testing.T) {
	handler := setup()

	t.Run("Error wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/v1/adverts", nil)
		handler.Mux.ServeHTTP(rec, req)
		wantResult := `{}`
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		} else if rec.Body.String() != wantResult {
			t.Fatalf("want: %v, got: %v", wantResult, rec.Body.String())
		}
	})

}

func TestParticularGroup(t *testing.T) {
	handler := setup()

	tests := []struct {
		name       string
		url        string
		method     string
		wantStatus int
		wantResult string
	}{
		{
			name:       "Error page not found",
			url:        "/v1/adverts/adc/1",
			method:     "GET",
			wantStatus: http.StatusNotFound,
			wantResult: `{}`,
		},
		{
			name:       "Error negative id",
			url:        "/v1/adverts/-5",
			method:     "GET",
			wantStatus: http.StatusNotFound,
			wantResult: `{}`,
		},
		{
			name:       "Error wrong method",
			url:        "/v1/adverts/5",
			method:     "POST",
			wantStatus: http.StatusMethodNotAllowed,
			wantResult: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			handler.Mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			} else if rec.Body.String() != tt.wantResult {
				t.Fatalf("want: %v, got: %v", tt.wantResult, rec.Body.String())
			}
		})
	}
}
