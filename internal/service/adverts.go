package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mrsubudei/adv-store-service/internal/entity"
	"github.com/mrsubudei/adv-store-service/internal/repository"
)

type AdvertService struct {
	repo repository.Advert
}

func NewAdvertService(repo repository.Advert) *AdvertService {
	return &AdvertService{
		repo: repo,
	}
}

func getTime() string {
	timeNow := time.Now()
	return timeNow.Format(DateFormat)
}

func (s *AdvertService) Create(ctx context.Context, adv entity.Advert) (int64, error) {
	adv.CreatedAt = getTime()

	err := s.repo.Store(ctx, &adv)
	if err != nil {
		if strings.Contains(err.Error(), UniqueNameConstraint) {
			return 0, entity.ErrNameAlreadyExist
		}
		return 0, fmt.Errorf("AdvertService - Create: %w", err)
	}

	return adv.Id, nil
}

func (s *AdvertService) GetById(ctx context.Context, id int64) (entity.Advert, error) {
	adv, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Advert{}, entity.ErrItemNotExists
		}
		return adv, fmt.Errorf("AdvertService - GetById: %w", err)
	}
	return adv, nil
}

func (s *AdvertService) GetAll(ctx context.Context) ([]entity.Advert, error) {
	adverts, err := s.repo.Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf("AdvertService - GetAll: %w", err)
	}
	if len(adverts) == 0 {
		return nil, entity.ErrNoItems
	}
	return adverts, nil
}

func (s *AdvertService) Update(ctx context.Context, adv entity.Advert) error {
	exist, err := s.repo.GetById(ctx, adv.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrItemNotExists
		}
		return fmt.Errorf("AdvertService - Update: %w", err)
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

	if adv.MainPhotoUrl != "" {
		exist.MainPhotoUrl = adv.MainPhotoUrl
	}

	if len(adv.PhotosUrls) != 0 {
		exist.PhotosUrls = []string{}
		exist.PhotosUrls = append(exist.PhotosUrls, adv.PhotosUrls...)
	}

	err = s.repo.Update(ctx, exist)
	if err != nil {
		if strings.Contains(err.Error(), UniqueNameConstraint) {
			return entity.ErrNameAlreadyExist
		}
		return fmt.Errorf("AdvertService - Update: %w", err)
	}

	return nil
}

func (s *AdvertService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrItemNotExists
		}
		return fmt.Errorf("AdvertService - Delete: %w", err)
	}

	return nil
}
