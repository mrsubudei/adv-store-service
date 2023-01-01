package v1_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrsubudei/adv-store-service/internal/config"
	v1 "github.com/mrsubudei/adv-store-service/internal/controller/http/v1"
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

func getMockHandler(t *testing.T, key string) http.HandlerFunc {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
	return mockHandler
}

func TestCheckAuth(t *testing.T) {
	// handler := setup()

	t.Run("OK", func(t *testing.T) {
		// if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		// 	t.Fatal(err)
		// }

		// mockHandler := getMockHandlerOne(t, "content")

		// handlerToTest := handler.CheckAuth(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}

		req.AddCookie(cookie)
		// handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

}
