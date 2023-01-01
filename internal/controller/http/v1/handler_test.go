package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
