package mock_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

type MockService struct {
	Adverts []entity.Advert
	Ids     int64
}

func NewMockService() *MockService {
	return &MockService{}
}

func (ms *MockService) Create(ctx context.Context, adv entity.Advert) (int64, error) {
	ms.Ids++
	adv.MainPhotoUrl = adv.PhotosUrls[0]
	adv.Id = ms.Ids

	for i := 0; i < len(ms.Adverts); i++ {
		if ms.Adverts[i].Name == adv.Name {
			return 0, entity.ErrNameAlreadyExist
		}
	}
	ms.Adverts = append(ms.Adverts, adv)

	return ms.Ids, nil
}

func (ms *MockService) GetById(ctx context.Context, id int64) (entity.Advert, error) {

	for i := 0; i < len(ms.Adverts); i++ {
		if ms.Adverts[i].Id == id {
			adv := ms.Adverts[i]
			adv.Id = 0
			return adv, nil
		}
	}
	return entity.Advert{}, entity.ErrItemNotExists
}

func (ms *MockService) getById(ctx context.Context, id int64) (*entity.Advert, error) {
	for i := 0; i < len(ms.Adverts); i++ {
		if ms.Adverts[i].Id == id {
			return &ms.Adverts[i], nil
		}
	}
	return &entity.Advert{}, entity.ErrItemNotExists
}

func (ms *MockService) GetAll(ctx context.Context) ([]entity.Advert, error) {
	if len(ms.Adverts) == 0 {
		return nil, entity.ErrNoItems
	}

	return ms.Adverts, nil
}

func (ms *MockService) Update(ctx context.Context, adv entity.Advert) error {
	exist, err := ms.getById(ctx, adv.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrItemNotExists
		}
		return fmt.Errorf("AdvertService - Update: %w", err)
	}
	fmt.Println(exist)
	fmt.Println(adv)
	if exist.Name == adv.Name {
		return entity.ErrNameAlreadyExist
	}

	if adv.Name != "" {
		exist.Name = adv.Name
	}

	if adv.Description != "" {
		exist.Description = adv.Description
	}

	if adv.Price != 0 {
		exist.Price = adv.Price
	}

	if len(adv.PhotosUrls) != 0 {
		exist.MainPhotoUrl = adv.PhotosUrls[0]
		exist.PhotosUrls = []string{}
		exist.PhotosUrls = append(exist.PhotosUrls, adv.PhotosUrls...)
	}

	return nil
}

func (ms *MockService) Delete(ctx context.Context, id int64) error {
	newAdverts := []entity.Advert{}
	found := false

	for i, v := range ms.Adverts {
		if v.Id == id {
			found = true
			newAdverts = deleteElement(ms.Adverts, i)
		}
	}
	if found {
		ms.Adverts = []entity.Advert{}
		ms.Adverts = append(ms.Adverts, newAdverts...)
		return nil
	} else {
		return entity.ErrItemNotExists
	}
}

func deleteElement[C any](sl []C, index int) []C {
	newSlice := []C{}
	newSlice = append(newSlice, sl[:index]...)
	newSlice = append(newSlice, sl[index+1:]...)
	return newSlice
}
