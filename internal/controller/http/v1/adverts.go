package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

func (h *Handler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	var adv entity.Advert
	err := h.parseJson(w, r, &adv)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - parseJson: %w", err))
		return
	}

	if errAns := h.checkData(adv); errAns.Error != "" {
		h.writeResponse(w, errAns)
		return
	}

	id, err := h.Service.Create(r.Context(), adv)
	if err != nil {
		if errors.Is(err, entity.ErrNameAlreadyExist) {
			h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
			h.writeResponse(w, ErrMessage{code: http.StatusConflict,
				Error: fmt.Sprintf(ItemNameExists, adv.Name)})
			return
		}
		h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
		h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError})
		return
	}

	ans := Response{
		Data: []entity.Advert{{Id: id}},
		code: http.StatusCreated,
	}
	h.writeResponse(w, ans)
}

func (h *Handler) GetAllAdverts(w http.ResponseWriter, r *http.Request) {
	advs, err := h.Service.GetAll(r.Context())
	if err != nil {
		if errors.Is(err, entity.ErrNoItems) {
			h.l.WriteLog(fmt.Errorf("v1 - GetAllAdverts - h.Service.GetAll: %w", err))
			h.writeResponse(w, Response{code: http.StatusOK, Data: []entity.Advert{}})
			return
		}
		h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - h.Service.Create: %w", err))
		h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError})
		return
	}

	maxPage := advs[0].MaxCount
	meta := &MetaData{MaxPage: maxPage}
	ans := Response{
		Data: advs,
		code: http.StatusOK,
		Meta: meta,
	}
	h.writeResponse(w, ans)
}

func (h *Handler) GetAdvert(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.KeyId).(int64)

	found, err := h.Service.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrItemNotExists) {
			h.l.WriteLog(fmt.Errorf("v1 - GetAdvert - h.Service.GetById: %w", err))
			h.writeResponse(w, ErrMessage{code: http.StatusNotFound,
				Error: NoContentFound + strconv.Itoa(int(id))})
			return
		}
		h.l.WriteLog(fmt.Errorf("v1 - GetAdvert - h.Service.Create: %w", err))
		h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError})
		return
	}

	ans := Response{code: http.StatusOK}

	if queryFields, ok := r.Context().Value(entity.KeyFields).(string); ok && queryFields != "" {
		ans.Data = []entity.Advert{found}
	} else {
		partialAdv := entity.Advert{Name: found.Name, Price: found.Price,
			MainPhotoUrl: found.MainPhotoUrl}
		ans.Data = []entity.Advert{partialAdv}
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

	id := r.Context().Value(entity.KeyId).(int64)
	adv.Id = id

	err = h.Service.Update(r.Context(), adv)
	if err != nil {
		if errors.Is(err, entity.ErrItemNotExists) {
			h.l.WriteLog(fmt.Errorf("v1 - UpdateAdvert - h.Service.Update #1: %w", err))
			h.writeResponse(w, ErrMessage{code: http.StatusNotFound,
				Error: NoContentFound + strconv.Itoa(int(id))})
			return
		} else if errors.Is(err, entity.ErrNameAlreadyExist) {
			h.l.WriteLog(fmt.Errorf("v1 - UpdateAdvert - h.Service.Update #2: %w", err))
			h.writeResponse(w, ErrMessage{code: http.StatusConflict,
				Error: fmt.Sprintf(ItemNameExists, adv.Name)})
			return
		}
		h.l.WriteLog(fmt.Errorf("v1 - UpdateAdvert - h.Service.Create: %w", err))
		h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError})
		return
	}

	ans := Response{
		code: http.StatusOK,
	}
	h.writeResponse(w, ans)
}

func (h *Handler) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.KeyId).(int64)

	err := h.Service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, entity.ErrItemNotExists) {
			h.l.WriteLog(fmt.Errorf("v1 - DeleteAdvert - h.Service.Delete: %w", err))
			h.writeResponse(w, ErrMessage{code: http.StatusNotFound,
				Error: NoContentFound + strconv.Itoa(int(id))})
			return
		}
		h.l.WriteLog(fmt.Errorf("v1 - DeleteAdvert - h.Service.Create: %w", err))
		h.writeResponse(w, ErrMessage{code: http.StatusInternalServerError})
		return
	}

	ans := Response{
		code: http.StatusOK,
	}
	h.writeResponse(w, ans)
}
