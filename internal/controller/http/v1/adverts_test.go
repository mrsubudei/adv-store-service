package v1_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

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
	longName := ""
	for i := 0; i < 201; i++ {
		longName += "a"
	}
	longDescription := ""
	for i := 0; i < 1001; i++ {
		longDescription += "a"
	}

	reqWithLongName := entity.Advert{Name: longName}
	jsonLongName, err := json.Marshal(reqWithLongName)
	if err != nil {
		t.Fatal(err)
	}
	reqWithLongDesc := entity.Advert{Name: "abc", Description: longDescription}
	jsonLongDesc, err := json.Marshal(reqWithLongDesc)
	if err != nil {
		t.Fatal(err)
	}

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
			name:       "Error too long name",
			reqData:    string(jsonLongName),
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"Request Entity Too Large","detail":"'name:' field's length exceeded"}`,
		},
		{
			name:       "Error too long description",
			reqData:    string(jsonLongDesc),
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"Request Entity Too Large","detail":"'description:' field's length exceeded"}`,
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
			name: "Error empty field: name",
			reqData: `{"description":"asd","price":40,
			"photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"request has empty fields","detail":"'name:' field is required"}`,
		},
		{
			name: "Error empty field: description",
			reqData: `{"name":"first item","price":40,
			"photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"request has empty fields","detail":"'description:' field is required"}`,
		},
		{
			name:       "Error empty field: price",
			reqData:    `{"name":"first item","description":"asd"}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"request has empty fields","detail":"'price:' field is required"}`,
		},
		{
			name:       "Error empty field: photo_urls",
			reqData:    `{"name":"first item","description":"asd","price":40}`,
			wantStatus: http.StatusBadRequest,
			wantResult: `{"error":"request has empty fields","detail":"'photo_urls:' field should have at least 1 url"}`,
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

	if _, err := handler.Service.Create(ctx, advert1); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.Service.Create(ctx, advert2); err != nil {
		t.Fatal(err)
	}

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
	ctx := context.Background()
	if _, err := handler.Service.Create(ctx, advert1); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.Service.Create(ctx, advert2); err != nil {
		t.Fatal(err)
	}

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
	if _, err := handler.Service.Create(context.Background(), advert1); err != nil {
		t.Fatal(err)
	}

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
		{
			name:       "Error item with that name already exists",
			wantStatus: http.StatusConflict,
			url:        "/v1/adverts/1",
			reqData:    `{"name":"new name","photo_urls":["http://files.com/12","http://files.com/13"]}`,
			wantResult: `{"error":"item with name 'new name' already exists"}`,
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
	if _, err := handler.Service.Create(context.Background(), advert1); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantResult string
	}{
		{
			name:       "OK",
			wantStatus: http.StatusNoContent,
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
