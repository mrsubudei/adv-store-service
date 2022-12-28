package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/mrsubudei/adv-store-service/internal/config"
	"github.com/mrsubudei/adv-store-service/internal/entity"
	"github.com/mrsubudei/adv-store-service/internal/service"
	"github.com/mrsubudei/adv-store-service/pkg/logger"
)

type Anwser struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
	code    int
}

type Handler struct {
	Service *service.AdvertService
	Cfg     config.Config
	l       *logger.Logger
	Mux     *http.ServeMux
}

func NewHandler(advService *service.AdvertService, cfg config.Config, logger *logger.Logger) *Handler {
	mux := http.NewServeMux()
	return &Handler{
		Service: advService,
		Cfg:     cfg,
		l:       logger,
		Mux:     mux,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/adverts/create", h.CreateAdvert)
}

func (h *Handler) parseJson(w http.ResponseWriter, r *http.Request, adv *entity.Advert) error {
	ans := Anwser{}
	wrongRequest := false

	err := json.NewDecoder(r.Body).Decode(adv)
	if err != nil {
		ans.code = http.StatusBadRequest
		ans.Message = WrongDataFormat
		wrongRequest = true
	}

	if utf8.RuneCountInString(adv.Description) > 1000 {
		ans.code = http.StatusBadRequest
		ans.Message = DescLengthExceeded
		wrongRequest = true
	} else if utf8.RuneCountInString(adv.Name) > 200 {
		ans.code = http.StatusBadRequest
		ans.Message = NameLengthExceeded
		wrongRequest = true
	} else if len(adv.PhotosUrls) > 3 {
		ans.code = http.StatusBadRequest
		ans.Message = UrlsNumberExceeded
		wrongRequest = true
	}

	if adv.Name == "" || adv.Description == "" || adv.Price == 0 || len(adv.PhotosUrls) == 0 {
		ans.code = http.StatusBadRequest
		ans.Message = EmptyFiledRequest
		wrongRequest = true
	}

	h.writeResponse(w, ans)

	if wrongRequest {
		return fmt.Errorf(WrongDataFormat)
	}
	return nil
}

func (h *Handler) writeResponse(w http.ResponseWriter, ans Anwser) {
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(ans)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - writeResponse - Marshal: %w", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(ans.code)
	w.Write(jsonResp)
}
