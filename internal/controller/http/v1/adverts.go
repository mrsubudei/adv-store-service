package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

func (h *Handler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	var adv entity.Advert
	err := h.parseJson(w, r, &adv)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - parseJson: %w", err))
		return
	}

	if errAns := h.checkData(adv); errAns.code != 0 {
		h.writeResponse(w, errAns)
		return
	}

	id, err := h.Service.Create(r.Context(), adv)
	if err != nil {
		if errors.Is(err, entity.ErrNameAlreadyExist) {
			h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
			h.writeResponse(w, ErrMessage{code: http.StatusConflict,
				Error:  http.StatusText(http.StatusConflict),
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
		code: http.StatusCreated,
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
	maxPage := advs[0].MaxCount
	ans := MultiResponse{
		Data: advs,
		code: http.StatusAccepted,
		Meta: MetaData{MaxPage: maxPage},
	}
	h.writeResponse(w, ans)
}

func (h *Handler) GetAdvert(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.KeyId).(int64)

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

	ans := SingleResponse{code: http.StatusAccepted}

	if queryFields := r.Context().Value(entity.KeyFields).(string); queryFields != "" {
		ans.Data = found
	} else {
		ans.Data.Name = found.Name
		ans.Data.Price = found.Price
		ans.Data.MainPhotoUrl = found.MainPhotoUrl
	}

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
				Error: "item not found"})
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
