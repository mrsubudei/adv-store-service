package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"encoding/json"

	"github.com/mrsubudei/adv-store-service/internal/config"
	"github.com/mrsubudei/adv-store-service/internal/entity"
	"github.com/mrsubudei/adv-store-service/internal/service"
	"github.com/mrsubudei/adv-store-service/pkg/logger"
)

type Handler struct {
	Service service.Service
	Cfg     config.Config
	l       *logger.Logger
	Mux     *http.ServeMux
}

func NewHandler(advService service.Service, cfg config.Config,
	logger *logger.Logger) *Handler {
	mux := http.NewServeMux()
	return &Handler{
		Service: advService,
		Cfg:     cfg,
		l:       logger,
		Mux:     mux,
	}
}

func (h *Handler) NewRouteGroups() {
	h.Mux.Handle("/v1/adverts", h.ParseQuery(http.HandlerFunc(h.CommonGroup)))
	h.Mux.Handle("/v1/adverts/", h.ParseQuery(http.HandlerFunc(h.ParticularGroup)))
	h.Mux.HandleFunc("/", h.WrongRoute)
}

func (h *Handler) WrongRoute(w http.ResponseWriter, r *http.Request) {
	h.writeResponse(w, ErrMessage{code: http.StatusNotFound})
}

func (h *Handler) CommonGroup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllAdverts(w, r)
	case http.MethodPost:
		h.CreateAdvert(w, r)
	default:
		h.writeResponse(w, ErrMessage{code: http.StatusMethodNotAllowed})
	}
}

func (h *Handler) ParticularGroup(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[len(path)-1])
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - NewPluralRoutes - Atoi: %w", err))
	}

	if id <= 0 || err != nil || r.URL.Path != "/v1/adverts/"+strconv.Itoa(id) {
		h.writeResponse(w, ErrMessage{code: http.StatusNotFound})
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
		h.writeResponse(w, ErrMessage{code: http.StatusMethodNotAllowed})
	}
}

func (h *Handler) parseJson(w http.ResponseWriter, r *http.Request, adv *entity.Advert) error {
	err := json.NewDecoder(r.Body).Decode(adv)
	if err != nil {
		h.writeResponse(w, ErrMessage{code: http.StatusBadRequest,
			Error: JsonNotCorrect})
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
	if _, err = w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) checkData(adv entity.Advert) ErrMessage {
	errMsg := ErrMessage{code: http.StatusBadRequest}
	switch {
	case utf8.RuneCountInString(adv.Description) > 1000:
		errMsg.Error = http.StatusText(http.StatusRequestEntityTooLarge)
		errMsg.Detail = DescLengthExceeded
	case utf8.RuneCountInString(adv.Name) > 200:
		errMsg.Error = http.StatusText(http.StatusRequestEntityTooLarge)
		errMsg.Detail = NameLengthExceeded
	case len(adv.PhotosUrls) > 3:
		errMsg.Error = http.StatusText(http.StatusRequestEntityTooLarge)
		errMsg.Detail = UrlsNumberExceeded
	case adv.Name == "":
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `'name:' field is required`
	case adv.Description == "":
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `'description:' field is required`
	case adv.Price == 0:
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `'price:' field is required`
	case len(adv.PhotosUrls) == 0:
		errMsg.Error = EmptyFiledRequest
		errMsg.Detail = `'photo_urls:' field should have at least 1 url`
	}

	return errMsg
}
