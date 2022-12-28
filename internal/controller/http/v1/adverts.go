package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

func (h *Handler) GetAdvert(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	h.Errors(w, http.StatusMethodNotAllowed)
	// 	return
	// }
	// if r.URL.Path != "/api/adverts" {
	// 	h.Errors(w, http.StatusNotFound)
	// 	return
	// }

	// adverts, err := h.Service.GetAll(r.Context())
	// if err != nil {
	// 	h.l.WriteLog(fmt.Errorf("v1 - GetAdvert - GetAllPosts - GetAll: %w", err))
	// 	h.Errors(w, http.StatusInternalServerError)
	// 	return
	// }

	// w.Write()
}

func (h *Handler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeResponse(w, Anwser{code: http.StatusMethodNotAllowed,
			Message: http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	var adv entity.Advert
	err := h.parseJson(w, r, &adv)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - parseJson: %w", err))
		return
	}

	id, err := h.Service.Create(r.Context(), adv)
	if err != nil {
		if errors.Is(err, entity.ErrNameAlreadyExist) {
			h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - Create: %w", err))
			h.writeResponse(w, Anwser{code: http.StatusConflict,
				Message: entity.ErrNameAlreadyExist.Error()})
			return
		}
		h.l.WriteLog(fmt.Errorf("v1 - CreateAdvert - Create: %w", err))
		h.writeResponse(w, Anwser{code: http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError)})
		return
	}

	ans := Anwser{id, AdvertCreated, http.StatusCreated}
	h.writeResponse(w, ans)
}
