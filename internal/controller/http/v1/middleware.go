package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

func (h *Handler) ParseQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//checking queries
		errMsg := ErrMessage{code: http.StatusBadRequest, Error: WrongQueryRequest}

		if val := r.URL.Query().Get(QueryFields); val != "" && val != QueryValueTrue {
			errMsg.Detail = `'fields=' query value should be 'true'`
		}
		if val := r.URL.Query().Get(QuerySortBy); val != "" && val != QueryValueCreatedAt &&
			val != QueryValuePrice {
			errMsg.Detail = `'sort_by=' query value should be either 'created_at' or 'price'`
		}
		if val := r.URL.Query().Get(QueryOrderBy); val != "" && val != QueryValueAsc &&
			val != QueryValueDesc {
			errMsg.Detail = `'order_by=' query value should be either 'asc' or 'desc'`
		}
		if val := r.URL.Query().Get(QueryOffset); val != "" {
			if parsedToInt, err := strconv.Atoi(val); err != nil || parsedToInt <= 0 {
				errMsg.Detail = `'offset=' query value should be positive number`
			}
		}
		if val := r.URL.Query().Get(QueryLimit); val != "" {
			if parsedToInt, err := strconv.Atoi(val); err != nil || parsedToInt <= 0 {
				errMsg.Detail = `'limit=' query value should be positive number`
			}
		}
		if errMsg.Detail != "" {
			h.writeResponse(w, errMsg)
			return
		}

		//parsing and adding queries to context
		queries := []string{QueryLimit, QueryOffset, QuerySortBy, QueryOrderBy, QueryFields}
		keys := []entity.ContextKey{entity.KeyLimit, entity.KeyOffset, entity.KeySortBy,
			entity.KeyOrderBy, entity.KeyFields}
		ctx := context.Background()
		for i := 0; i < len(queries); i++ {
			if value := r.URL.Query().Get(queries[i]); value != "" {
				if parsedToInt, err := strconv.Atoi(value); err == nil {
					ctx = context.WithValue(ctx, keys[i], parsedToInt)
				} else {
					ctx = context.WithValue(ctx, keys[i], value)
				}
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
