package v1_test

// func TestCreateAdvert(t *testing.T) {
// 	handler := setup()

// 	t.Run("OK", func(t *testing.T) {
// 		rec := httptest.NewRecorder()
// 		req := httptest.NewRequest(http.MethodGet, "/v1/adverts", nil)
// 		handler.Mux.ServeHTTP(rec, req)

// 		if rec.Code != http.StatusOK {
// 			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
// 		}
// 	})

// 	t.Run("err wrong method", func(t *testing.T) {
// 		rec := httptest.NewRecorder()
// 		req := httptest.NewRequest(http.MethodPut, "/", nil)
// 		handler.Mux.ServeHTTP(rec, req)

// 		if rec.Code != http.StatusMethodNotAllowed {
// 			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
// 		}
// 	})
// }
