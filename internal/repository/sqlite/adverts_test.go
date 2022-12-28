package sqlite_test

import (
        "reflect"
        "testing"
        "context"
        "strings"
        "errors"
        "database/sql"

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
    advert3 = entity.Advert{
        Name: "suit",
        Description: "Lorem ipsum dolor",
        Price: 50,
        MainPhotoUrl: "http:fs.com/7",
        PhotosUrls: []string{
          "http:fs.com/8",
          "http:fs.com/9",
          "http:fs.com/10",
        },
    }   
)

func TestStore(t *testing.T) {
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

        })

        t.Run("Err name already exist", func(t *testing.T) {
            advRepeatedName := entity.Advert{
                Name: "car",
            }
        
            if err = repo.Store(ctx, &advRepeatedName); err == nil {
                        t.Fatalf("Error expected")
            } else if !strings.Contains(err.Error(), "UNIQUE constraint failed: adverts.name") {
                        t.Fatalf("unexpected error: %v", err)
            }
        })
}

func TestGetById(t *testing.T) {
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
                
                if err := repo.Store(ctx, &advert2); err != nil {
                        t.Fatal("Unable to store:", err)
                }

                if found, err := repo.GetById(ctx, 1); err != nil {
                        t.Fatal("Unable to GetById:", err)
                } else if !reflect.DeepEqual(advert1, found) {
                        t.Fatalf("mismatch: %#v != %#v", advert1, found)
                }
                
                if found, err := repo.GetById(ctx, 2); err != nil {
                        t.Fatal("Unable to GetById:", err)
                } else if !reflect.DeepEqual(advert2, found) {
                        t.Fatalf("mismatch: %#v != %#v", advert2, found)
                }
        })

        t.Run("Err did not found", func(t *testing.T) {
                if _, err := repo.GetById(ctx, 4); err == nil {
                         t.Fatalf("Error expected")
                } else if !errors.Is(err, sql.ErrNoRows) {
                        t.Fatalf("want: %v, got: %v", sql.ErrNoRows, err)
                }
        })
}

func TestFetch(t *testing.T) {
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
                
                if err := repo.Store(ctx, &advert2); err != nil {
                        t.Fatal("Unable to store:", err)
                }
                
                 if err := repo.Store(ctx, &advert3); err != nil {
                        t.Fatal("Unable to store:", err)
                }

                if found, err := repo.Fetch(ctx); err != nil {
                        t.Fatal("Unable to Fetch:", err)
                } else if len(found) != 3 {
                        t.Fatalf("want: %d, got: %d", 3, len(found))
                } 
        })        
}

func TestUpdate(t *testing.T) {
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
                
                updatedAdv := entity.Advert{
                    Id: 1,
                    Name: "updated name",
                    Description: "updated desc",
                    Price: 777,
                    MainPhotoUrl: "new url 1",
                    PhotosUrls: []string{
                       "new url 1",
                       "new url 2",
                       "new url 3",
                    },
                }   
                
                if err := repo.Update(ctx, updatedAdv); err != nil {
                    t.Fatal("Unable to Fetch:", err)
                } 
                
                if found, err := repo.GetById(ctx, 1); err != nil {
                        t.Fatal("Unable to GetById:", err)
                } else if !reflect.DeepEqual(updatedAdv, found) {
                        t.Fatalf("mismatch: %#v != %#v", updatedAdv, found)
                }
        })        
}

func TestDelete(t *testing.T) {
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
                
                if err := repo.Delete(ctx, 1); err != nil {
                        t.Fatal("Unable to GetById:", err)
                } 
                
                if _, err := repo.GetById(ctx, 1); err == nil {
                        t.Fatal("Error expected")
                } else if !errors.Is(err, sql.ErrNoRows) {
                        t.Fatalf("want: %v, got: %v", sql.ErrNoRows, err)
                }
        })        
}
