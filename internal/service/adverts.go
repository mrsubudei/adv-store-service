package service

import (
	"fmt"
  "context"
	"github.com/mrsubudei/adv-store-service/internal/entity"
  "github.com/mrsubudei/adv-store-service/internal/repository/sqlite"
)

type Service struct {
    repo sqlite.AdvertsRepo
}

func NewService(repo sqlite.AdvertsRepo) *Service {
  return &Service{
    repo: repo,
  } 
}

func (s *Service) Create(ctx context.Context, adv entity.Advert) error {
    err := s.repo.Store(ctx, &adv)
    if err != nil {
        return fmt.Errorf("Service - Create: %w", err) 
    }
  
  return nil
}

func (s *Service) GetById(ctx context.Context, id int64) (entity.Advert, error) {
    advert := entity.Advert{}
    adv, err := s.repo.GetById(ctx, id)
     if err != nil {
        return adv, fmt.Errorf("Service - GetById: %w", err) 
     }
     return advert, nil
}

func (s *Service) GetAll(ctx context.Context) ([]entity.Advert, error) {
    adverts, err := s.repo.Fetch(ctx)
    if err != nil {
        return nil, fmt.Errorf("Service - GetAll: %w", err) 
     }
     return adverts, nil
}

func (s *Service) Update(ctx context.Context, adv entity.Advert) error {
    exist, err := s.repo.GetById(ctx, adv.Id)
    if err != nil {
        return fmt.Errorf("Service - Update: %w", err) 
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
        return fmt.Errorf("Service - Update: %w", err) 
    }
    
    return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
    err := s.repo.Delete(ctx, id)
    if err != nil {
        return fmt.Errorf("Service - Delete: %w", err) 
    }
    
    return nil
}


