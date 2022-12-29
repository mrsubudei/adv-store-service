package v1

import (
        "encoding/json"
        "fmt"
        "net/http"
        "strings"
       "strconv"
       "context"

        "github.com/mrsubudei/adv-store-service/internal/config"
        "github.com/mrsubudei/adv-store-service/internal/entity"
        "github.com/mrsubudei/adv-store-service/internal/service"
        "github.com/mrsubudei/adv-store-service/pkg/logger"
)

type SingleResponse struct {
    Data entity.Advert `json:"data"`
    code    int
}

type MultiResponse struct {
        Data    []entity.Advert `json:"data"`
        code    int
}

type Handler struct {
        Service *service.AdvertService
        Cfg     config.Config
        l       *logger.Logger
        Mux     *http.ServeMux
}

func (r SingleResponse) getCode() int{
    return r.code
}

func (r MultiResponse) getCode() int{
    return r.code
}

func (e ErrMessage) getCode() int{
    return e.code
}

type Answer interface {
    getCode() int
}

type ContextKey string

const KeyId ContextKey = "id"

func NewHandler(advService *service.AdvertService, cfg config.Config, logger *logger.Logger) *Handler {
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
        fmt.Println(id)
    if len(path) > 3 || len(path) == 1 || path[1] != "adverts" || id <= 0 || err != nil {
        h.writeResponse(w, ErrMessage{Error: http.StatusText(http.StatusNotFound), 
        code: http.StatusNotFound})
        return
    }

        ctx := context.WithValue(r.Context(), KeyId, id)

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
