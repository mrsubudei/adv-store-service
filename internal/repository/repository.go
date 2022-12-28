package repository

import (
  "context" 
  "github.com/mrsubudei/adv-store-service/internal/entity"
)

type Advert interface{
  Store(ctx context.Context, adv *entity.Advert) error
  GetById(ctx context.Context, id int64) (entity.Advert, error)
  Fetch(ctx context.Context) ([]entity.Advert, error)
  Update(ctx context.Context, adv entity.Advert) error 
  Delete(ctx context.Context, id int64) error
}
