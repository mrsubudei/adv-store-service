package v1

import (
        "errors"
        "fmt"
        "net/http"
        "unicode/utf8"

        "github.com/mrsubudei/adv-store-service/internal/entity"
)

func (h *Handler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
        var adv entity.Advert
        err := h.parseJson(w, r, &adv)
        if err != nil {
                h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - parseJson: %w", err))
                return
        }
        
        wrongRequest := false
        errAns := ErrMessage{}
        
        switch {
            case utf8.RuneCountInString(adv.Description) > 1000:
                errAns.code = http.StatusBadRequest
                errAns.Error = http.StatusText(http.StatusRequestEntityTooLarge)
                errAns.Detail = DescLengthExceeded
                wrongRequest = true
            case utf8.RuneCountInString(adv.Name) > 200:
                errAns.code = http.StatusBadRequest
                errAns.Error = http.StatusText(http.StatusRequestEntityTooLarge)
                errAns.Detail = NameLengthExceeded
                wrongRequest = true
            case len(adv.PhotosUrls) > 3:
                errAns.code = http.StatusBadRequest
                errAns.Error = http.StatusText(http.StatusRequestEntityTooLarge)
                errAns.Detail = UrlsNumberExceeded
                wrongRequest = true
            case adv.Name == "":
                errAns.code = http.StatusBadRequest
                errAns.Error = EmptyFiledRequest
                errAns.Detail = "\"name\": field is required"
                wrongRequest = true
            case adv.Description == "":
                errAns.code = http.StatusBadRequest
                errAns.Error = EmptyFiledRequest
                errAns.Detail = "\"description\": field is required"
                wrongRequest = true
            case adv.Price == 0:
                errAns.code = http.StatusBadRequest
                errAns.Error = EmptyFiledRequest
                errAns.Detail = "\"price\": field is required"
                wrongRequest = true
            case len(adv.PhotosUrls) == 0:
                errAns.code = http.StatusBadRequest
                errAns.Error = EmptyFiledRequest
                errAns.Detail = "\"photo_urls\": field should have at least 1 url"
                wrongRequest = true
        }

        if wrongRequest {
            h.writeResponse(w, errAns)
            return
        }

        id, err := h.Service.Create(r.Context(), adv)
        if err != nil {
                if errors.Is(err, entity.ErrNameAlreadyExist) {
                        h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
                        h.writeResponse(w, ErrMessage{code: http.StatusConflict, Error: http.StatusText(http.StatusConflict),
                                Detail: entity.ErrNameAlreadyExist.Error()})
                        return
                }
                h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
                h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError,
                        Error: http.StatusText(http.StatusInternalServerError), Detail: err.Error()})
                return
        }
    
        ans := SingleResponse{
          Data: entity.Advert{Id: id},       
          code:   http.StatusCreated,
        }
        h.writeResponse(w, ans)
}

func (h *Handler) GetAllAdverts(w http.ResponseWriter, r *http.Request) {
    advs, err := h.Service.GetAll(r.Context())
    if err != nil {
        if errors.Is(err, entity.ErrNoItems) {
            h.l.WriteLog(fmt.Errorf("v1 - GetAllAdverts - h.Service.GetAll: %w", err))
             h.writeResponse(w, ErrMessage{code: http.StatusNoContent, 
             Error: http.StatusText(http.StatusNoContent)})
            return
        }
        h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
        h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError,
            Error: http.StatusText(http.StatusInternalServerError), Detail: err.Error()})
        return
    }
    
    ans := MultiResponse{
        Data: advs,
        code: http.StatusAccepted,
    }
     h.writeResponse(w, ans)
}

func (h *Handler) GetAdvert(w http.ResponseWriter, r *http.Request) {
        id, ok := r.Context().Value(KeyId).(int64)
        if !ok {
            h.l.WriteLog(fmt.Errorf("v1 - GetAdvert - TypeAssertion:"+
                        "got data of type %T but wanted int", id))
            h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError,
            Error: http.StatusText(http.StatusInternalServerError)})
            return
        }
        fmt.Println(id)        
        found, err := h.Service.GetById(r.Context(), id)
        if err != nil {
                if errors.Is(err, entity.ErrItemNotExists) {
                        h.l.WriteLog(fmt.Errorf("v1 - GetAdvert - h.Service.GetById: %w", err))
                        h.writeResponse(w, ErrMessage{code: http.StatusNoContent, 
                        Error: http.StatusText(http.StatusNoContent)})
                        return
                }
                h.l.WriteLog(fmt.Errorf("v1 - GetAdvert - h.Service.Create: %w", err))
                h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError,
                        Error: http.StatusText(http.StatusInternalServerError), Detail: err.Error()})
                return
        }
        
        ans := SingleResponse{code: http.StatusCreated}
        
        if len(r.URL.Query().Get("fields")) == 0 {
            ans.Data.Name = found.Name
            ans.Data.Price = found.Price
            ans.Data.MainPhotoUrl = found.MainPhotoUrl
        }
    
        ans.Data = found
        h.writeResponse(w, ans)
}

func (h *Handler) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
        var adv entity.Advert
        err := h.parseJson(w, r, &adv)
        if err != nil {
                h.l.WriteLog(fmt.Errorf("v1 - UpdateAdvert - parseJson: %w", err))
                return
        }
        
        err = h.Service.Update(r.Context(), adv)
        if err != nil {
                if errors.Is(err, entity.ErrItemNotExists) {
                        h.l.WriteLog(fmt.Errorf("v1 - UpdateAdvert - h.Service.Update: %w", err))
                        h.writeResponse(w, ErrMessage{code: http.StatusNoContent, 
                        Error: http.StatusText(http.StatusNoContent)})
                        return
                }
                h.l.WriteLog(fmt.Errorf("v1 - UpdateAdvert - h.Service.Create: %w", err))
                h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError,
                        Error: http.StatusText(http.StatusInternalServerError), Detail: err.Error()})
                return
        }
        
        ans := SingleResponse{
        code: http.StatusAccepted,
        }
        h.writeResponse(w, ans)
}

func (h *Handler) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
        var adv entity.Advert
        err := h.parseJson(w, r, &adv)
        if err != nil {
                h.l.WriteLog(fmt.Errorf("v1 - DeleteAdvert - parseJson: %w", err))
                return
        }
        
        err = h.Service.Delete(r.Context(), adv.Id)
        if err != nil {
                if errors.Is(err, entity.ErrItemNotExists) {
                        h.l.WriteLog(fmt.Errorf("v1 - DeleteAdvert - h.Service.Delete: %w", err))
                        h.writeResponse(w, ErrMessage{code: http.StatusNoContent, 
                        Error: http.StatusText(http.StatusNoContent)})
                        return
                }
                h.l.WriteLog(fmt.Errorf("v1 - DeleteAdvert - h.Service.Create: %w", err))
                h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError,
                        Error: http.StatusText(http.StatusInternalServerError), Detail: err.Error()})
                return
        }
        
        ans := SingleResponse{
        code: http.StatusAccepted,
        }
        h.writeResponse(w, ans)
}
