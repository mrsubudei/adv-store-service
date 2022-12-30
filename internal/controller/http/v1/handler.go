package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mrsubudei/adv-store-service/internal/config"
	"github.com/mrsubudei/adv-store-service/internal/entity"
	"github.com/mrsubudei/adv-store-service/internal/service"
	"github.com/mrsubudei/adv-store-service/pkg/logger"
)

type Handler struct {
	Service *service.AdvertService
	Cfg     config.Config
	l       *logger.Logger
	Mux     *http.ServeMux
}

func NewHandler(advService *service.AdvertService, cfg config.Config,
	logger *logger.Logger) *Handler {
	mux := http.NewServeMux()
	return &Handler{
		Service: advService,
		Cfg:     cfg,
		l:       logger,
		Mux:     mux,
	}
}

func (h *Handler) NewCommonRoutes(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
	case http.MethodGet:
		h.GetAllAdverts(w, r)
	case http.MethodPost:
		h.CreateAdvert(w, r)
	default:
		h.writeResponse(w, ErrMessage{Error: http.StatusText(http.StatusMethodNotAllowed),
			code: http.StatusMethodNotAllowed})
	}
}

func (h *Handler) NewParticularRoutes(w http.ResponseWriter, r *http.Request) {
    path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - NewPluralRoutes - Atoi: %w", err))
	}

	if len(path) > 4 || len(path) == 2 || path[2] != "adverts" || id <= 0 || err != nil {
		h.writeResponse(w, ErrMessage{Error: http.StatusText(http.StatusNotFound),
			code: http.StatusNotFound})
		return
	}

	ctx := context.WithValue(r.Context(), entity.KeyId, int64(id))

	switch r.Method {
	case http.MethodGet:
		h.GetAdvert(w, r.WithContext(ctx))
	case http.MethodPut:
		h.UpdateAdvert(w, r.WithContext(ctx))
	case http.MethodDelete:
		h.DeleteAdvert(w, r.WithContext(ctx))
	default:
		h.writeResponse(w, ErrMessage{Error: http.StatusText(http.StatusMethodNotAllowed),
			code: http.StatusMethodNotAllowed})
	}
}

func (h *Handler) parseJson(w http.ResponseWriter, r *http.Request, adv *entity.Advert) error {
	err := json.NewDecoder(r.Body).Decode(adv)
	if err != nil {
		h.writeResponse(w, ErrMessage{code: http.StatusBadRequest,
			Error: http.StatusText(http.StatusBadRequest), Detail: "format not correct"})
		return fmt.Errorf(WrongDataFormat)
	}

	return nil
}

func (h *Handler) writeResponse(w http.ResponseWriter, ans Answer) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	jsonResp, err := json.Marshal(ans)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - writeResponse - Marshal: %w", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(ans.getCode())
	w.Write(jsonResp)
}

func (h *Handler) checkData(adv entity.Advert) ErrMessage {
	errMsg := ErrMessage{}
	switch {
	case utf8.RuneCountInString(adv.Description) > 1000:
		errMsg.code = http.StatusBadRequest
		errMsg.Error = http.StatusText(http.StatusRequestEntityTooLarge)
		errMsg.Detail = DescLengthExceeded
	case utf8.RuneCountInString(adv.Name) > 200:
		errMsg.code = http.StatusBadRequest
		errMsg.Error = http.StatusText(http.StatusRequestEntityTooLarge)
		errMsg.Detail = NameLengthExceeded
	case len(adv.PhotosUrls) > 3:
		errMsg.code = http.StatusBadRequest
		errMsg.Error = http.StatusText(http.StatusRequestEntityTooLarge)
		errMsg.Detail = UrlsNumberExceeded
	case adv.Name == "":
		errMsg.code = http.StatusBadRequest
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `"name": field is required"`
	case adv.Description == "":
		errMsg.code = http.StatusBadRequest
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `"description": field is required"`
	case adv.Price == 0:
		errMsg.code = http.StatusBadRequest
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `"price": field is required`
	case len(adv.PhotosUrls) == 0:
		errMsg.code = http.StatusBadRequest
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `"photo_urls:" field should have at least 1 url`
	}
	
	return errMsg
}
 
func (h *Handler) CheckAndParseQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    //checking queries
	    errMsg := ErrMessage{code: http.StatusBadRequest, Error: WrongQueryRequest}
	
	if val := r.URL.Query().Get(QueryFields); val != "" && val != QueryValueTrue {
		errMsg.Detail = `"fields=" query value should be "true"`
	}
	if val := r.URL.Query().Get(QuerySortBy); val != "" && val != QueryValueCreatedAt && val != QueryValuePrice {
		errMsg.Detail = `"sort_by=" query value should be either "created_at" or "price"`
	}
	if val := r.URL.Query().Get(QueryOrderBy); val != "" && val != QueryValueAsc && val != QueryValueDesc {
		errMsg.Detail = `"order_by=" query value should be either "asc" or "desc"`
	}
	if val := r.URL.Query().Get(QueryOffset); val != "" {
	    if parsedToInt, err := strconv.Atoi(val); err != nil || parsedToInt <= 0 {
			errMsg.Detail = `"offset=" query value should be positive number`
		}	
	}
	if val := r.URL.Query().Get(QueryLimit); val != "" {
	    if parsedToInt, err := strconv.Atoi(val); err != nil || parsedToInt <= 0 {
			errMsg.Detail = `"limit=" query value should be positive number`
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
		      fmt.Println(parsedToInt)
			  ctx = context.WithValue(ctx, keys[i], parsedToInt)
		   } else {
		       ctx = context.WithValue(ctx, keys[i], value)
		   }
		}
	}

	 next.ServeHTTP(w, r.WithContext(ctx))
	})
}    
    
