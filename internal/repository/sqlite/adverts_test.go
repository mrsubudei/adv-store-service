package sqlite_test

import (
        "reflect"
        "strings"
        "testing"
        "context"

        "github.com/mrsubudei/adv-store-service/internal/entity"
        "github.com/mrsubudei/adv-store-service/internal/repository/sqlite"
)

var (
    advert1 = entity.Advert{
        Name: "car",
        Description: "Lorem ipsum dolor sit amet",
        Price: 150,
        MainPhotoUrl: "http:fs.com/1",
        PhotosUrls: []string{
          "http:fs.com/1",
          "http:fs.com/2",
          "http:fs.com/3",
        },
    }   
    
    advert2 = entity.Advert{
        Name: "toy",
        Description: "Lorem ipsum dolor sit",
        Price: 90,
        MainPhotoUrl: "http:fs.com/6",
        PhotosUrls: []string{
          "http:fs.com/6",
          "http:fs.com/7",
          "http:fs.com/8",
        },
    }   
)

func TestUserStore(t *testing.T) {
        db := sqlite.MustOpenDB(t, "file:foobar?mode=memory&cache=shared")
        defer sqlite.MustCloseDB(t, db)
        err := sqlite.CreateDB(db)
        if err != nil {
                t.Fatal("Unable to create db:", err)
        }
        repo := sqlite.NewAdvertsRepo(db)
        ctx := context.Background()
        
        t.Run("OK", func(t *testing.T) {
        
                if err := repo.Store(ctx, &advert1); err != nil {
                        t.Fatal("Unable to store:", err)
                }

                if found, err := repo.GetById(ctx, 1); err != nil {
                        t.Fatal("Unable to GetById:", err)
                } else if !reflect.DeepEqual(advert1, found) {
                        t.Fatalf("mismatch: %#v != %#v", advert1, found)
                }

                if err := repo.Store(ctx, &advert2); err != nil {
                        t.Fatal("Unable to Store:", err)
                }

                if found2, err := repo.GetById(ctx, 2); err != nil {
                        t.Fatal("Unable to GetById:", err)
                } else if !reflect.DeepEqual(advert2, found2) {
                        t.Fatalf("mismatch: %#v != %#v", advert2, found2)
                }
        })

        t.Run("ErrNameAlreadyExist", func(t *testing.T) {
                advert3 := entity.Advert{
            Name: "toy",
            Description: "Lorem ipsum dolor sit",
            Price: 90,
            MainPhotoUrl: "http:fs.com/6",
            PhotosUrls: []string{
                "http:fs.com/6",
                "http:fs.com/7",
                "http:fs.com/8",
            },
        }   
                if err = repo.Store(ctx, &advert3); err == nil {
                        t.Fatalf("Error expected")
                } else if !strings.Contains(err.Error(), "UNIQUE constraint failed: adverts.name") {
                        t.Fatalf("unexpected error: %v", err)
                }
        })
}
