package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mrsubudei/adv-store-service/internal/entity"
	m "github.com/mrsubudei/adv-store-service/internal/repository/mock"
	"github.com/mrsubudei/adv-store-service/internal/service"
)

var (
	advert1 = entity.Advert{
		Id:           1,
		Name:         "car",
		Description:  "Lorem ipsum dolor sit amet",
		Price:        150,
		MainPhotoUrl: "http:fs.com/1",
		PhotosUrls: []string{
			"http:fs.com/1",
			"http:fs.com/2",
			"http:fs.com/3",
		},
	}

	advert2 = entity.Advert{
		Id:           2,
		Name:         "toy",
		Description:  "Lorem ipsum dolor sit",
		Price:        90,
		MainPhotoUrl: "http:fs.com/6",
		PhotosUrls: []string{
			"http:fs.com/6",
			"http:fs.com/7",
			"http:fs.com/8",
		},
	}
	advert3 = entity.Advert{
		Id:           3,
		Name:         "suit",
		Description:  "Lorem ipsum dolor",
		Price:        50,
		MainPhotoUrl: "http:fs.com/7",
		PhotosUrls: []string{
			"http:fs.com/8",
			"http:fs.com/9",
			"http:fs.com/10",
		},
	}
)

func TestCreate(t *testing.T) {
	mockRepo := m.NewMockRepo()
	service := service.NewAdvertService(mockRepo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		if id, err := service.Create(ctx, advert1); err != nil {
			t.Fatal(err)
		} else if id != 1 {
			t.Fatalf("want: %d, got: %d", 1, id)
		}
	})

	t.Run("err name already exist", func(t *testing.T) {
		if _, err := service.Create(ctx, advert1); err == nil {
			t.Fatal("Expected error")
		} else if !errors.Is(err, entity.ErrNameAlreadyExist) {
			t.Fatalf("want: %v, got: %v", entity.ErrNameAlreadyExist, err)
		}
	})
}

func TestGetById(t *testing.T) {
	mockRepo := m.NewMockRepo()
	service := service.NewAdvertService(mockRepo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		var id int64
		var err error
		if id, err = service.Create(ctx, advert1); err != nil {
			t.Fatal(err)
		}
		found, err := service.GetById(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		advert1.CreatedAt = found.CreatedAt

		if !reflect.DeepEqual(advert1, found) {
			t.Fatalf("mismatch: %#v != %#v", advert1, found)
		}
	})

	t.Run("err item not found", func(t *testing.T) {
		if _, err := service.GetById(ctx, 89); err == nil {
			t.Fatal("Error expected")
		} else if !errors.Is(err, entity.ErrItemNotExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrItemNotExists, err)
		}
	})
}

func TestGetAll(t *testing.T) {
	mockRepo := m.NewMockRepo()
	service := service.NewAdvertService(mockRepo)
	ctx := context.Background()

	t.Run("Error no items", func(t *testing.T) {
		if _, err := service.GetAll(ctx); err == nil {
			t.Fatal("Error expected")
		} else if !errors.Is(err, entity.ErrNoItems) {
			t.Fatalf("want: %v, got: %v", entity.ErrNoItems, err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		if _, err := service.Create(ctx, advert1); err != nil {
			t.Fatal(err)
		}

		if _, err := service.Create(ctx, advert2); err != nil {
			t.Fatal(err)
		}

		if _, err := service.Create(ctx, advert3); err != nil {
			t.Fatal(err)
		}

		if found, err := service.GetAll(ctx); err != nil {
			t.Fatal(err)
		} else if len(found) != 3 {
			t.Fatalf("want: %d, got: %d", 3, len(found))
		}
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := m.NewMockRepo()
	service := service.NewAdvertService(mockRepo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		var id int64
		var err error
		if id, err = service.Create(ctx, advert1); err != nil {
			t.Fatal(err)
		}

		updated := entity.Advert{
			Id:           id,
			Name:         "updated name",
			Description:  "updated desc",
			Price:        777,
			MainPhotoUrl: "updated url 1",
			PhotosUrls: []string{
				"updated url 1",
				"updated url 2",
				"updated url 3",
			},
		}

		if err := service.Update(ctx, updated); err != nil {
			t.Fatal(err)
		}

		found, err := service.GetById(ctx, id)
		if err != nil {
			t.Fatal(err)
		}

		updated.CreatedAt = found.CreatedAt

		if !reflect.DeepEqual(updated, found) {
			t.Fatalf("mismatch: %#v != %#v", updated, found)
		}
	})

	t.Run("Err item not found", func(t *testing.T) {
		updated := entity.Advert{
			Id:           987,
			Name:         "updated name",
			Description:  "updated desc",
			Price:        777,
			MainPhotoUrl: "updated url 1",
			PhotosUrls: []string{
				"updated url 1",
				"updated url 2",
				"updated url 3",
			},
		}

		if err := service.Update(ctx, updated); err == nil {
			t.Fatal("Error expected")
		} else if !errors.Is(err, entity.ErrItemNotExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrItemNotExists, err)
		}
	})
}

func TestDelete(t *testing.T) {
	mockRepo := m.NewMockRepo()
	service := service.NewAdvertService(mockRepo)
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		var id int64
		var err error
		if id, err = service.Create(ctx, advert1); err != nil {
			t.Fatal(err)
		}

		if err := service.Delete(ctx, id); err != nil {
			t.Fatal(err)
		}

		if _, err := service.GetById(ctx, id); err == nil {
			t.Fatal("Error expected")
		} else if !errors.Is(err, entity.ErrItemNotExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrItemNotExists, err)
		}
	})

	t.Run("Err item not found", func(t *testing.T) {

		if err := service.Delete(ctx, 1); err == nil {
			t.Fatal("Error expected")
		} else if !errors.Is(err, entity.ErrItemNotExists) {
			t.Fatalf("want: %v, got: %v", entity.ErrItemNotExists, err)
		}
	})
}
