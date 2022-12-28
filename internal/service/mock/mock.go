package mock_service

import (
    "context"
     "github.com/mrsubudei/adv-store-service/internal/entity"
    )
    
type MockService struct {
    Adverts []entity.Advert
}

func NewMockService() *MockService {
    return &MockService{}
}

func (ms *MockService) Create(ctx context.Context, adv entity.Advert) error {
    for i:=0; i<len(ms.Adverts); i++ {
         if ms.Adverts[i].Id == adv.Id {
             return entity.ErrNameAlreadyExist
         }
     }
     ms.Adverts = append(ms.Adverts, adv)
     
     return nil
}

 func (ms *MockService) GetById(ctx context.Context, id int64) (entity.Advert, error) {
     for i:=0; i<len(ms.Adverts); i++ {
         if ms.Adverts[i].Id == id {
             return ms.Adverts[i], nil
         }
     }
     return entity.Advert{}, entity.ErrItemNotExists
 }
 
func (ms *MockService)  GetAll(ctx context.Context) ([]entity.Advert, error) {
    if len(ms.Adverts) == 0 {
        return nil, entity.ErrNoItems
    }
     return ms.Adverts, nil
}
 func (ms *MockService) Update(ctx context.Context, adv entity.Advert) error {
     for i:=0; i<len(ms.Adverts); i++ {
         if ms.Adverts[i].Id == adv.Id {
             ms.Adverts[i] = adv
             return nil
         }
     }
     return entity.ErrItemNotExists
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
 
 
