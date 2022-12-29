package mock_repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrsubudei/adv-store-service/internal/entity"
	"github.com/mrsubudei/adv-store-service/internal/service"
)

type MockRepo struct {
	Adverts []entity.Advert
}

func NewMockRepo() *MockRepo {
	return &MockRepo{}
}

func (mr *MockRepo) Store(ctx context.Context, adv *entity.Advert) error {
	for i := 0; i < len(mr.Adverts); i++ {
		if mr.Adverts[i].Id == adv.Id {
			return fmt.Errorf(service.UniqueNameConstraint)
		}
	}
	mr.Adverts = append(mr.Adverts, *adv)
	return nil
}
func (mr *MockRepo) GetById(ctx context.Context, id int64) (entity.Advert, error) {
	for i := 0; i < len(mr.Adverts); i++ {
		if mr.Adverts[i].Id == id {
			return mr.Adverts[i], nil
		}
	}

	return entity.Advert{}, sql.ErrNoRows
}
func (mr *MockRepo) Fetch(ctx context.Context) ([]entity.Advert, error) {
	return mr.Adverts, nil
}
func (mr *MockRepo) Update(ctx context.Context, adv entity.Advert) error {
	for i := 0; i < len(mr.Adverts); i++ {
		if mr.Adverts[i].Id == adv.Id {
			mr.Adverts[i] = adv
			return nil
		}

	}
	return sql.ErrNoRows
}

func (mr *MockRepo) Delete(ctx context.Context, id int64) error {
	newAdverts := []entity.Advert{}
	found := false
	for i, v := range mr.Adverts {
		if v.Id == id {
			found = true
			newAdverts = deleteElement(mr.Adverts, i)
		}
	}
	if found {
		mr.Adverts = []entity.Advert{}
		mr.Adverts = append(mr.Adverts, newAdverts...)
		return nil
	} else {
		return sql.ErrNoRows
	}
}

func deleteElement[C any](sl []C, index int) []C {
	newSlice := []C{}
	newSlice = append(newSlice, sl[:index]...)
	newSlice = append(newSlice, sl[index+1:]...)
	return newSlice
}
