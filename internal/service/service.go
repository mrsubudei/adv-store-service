package service

import (
	"context"

	"github.com/mrsubudei/adv-store-service/internal/entity"
)

type Service interface {
	Create(ctx context.Context, adv entity.Advert) (int64, error)
	GetById(ctx context.Context, id int64) (entity.Advert, error)
	GetAll(ctx context.Context) ([]entity.Advert, error)
	Update(ctx context.Context, adv entity.Advert) error
	Delete(ctx context.Context, id int64) error
}
