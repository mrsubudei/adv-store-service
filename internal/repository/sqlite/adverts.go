package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrsubudei/adv-store-service/internal/entity"
	"github.com/mrsubudei/adv-store-service/pkg/sqlite3"
)

type AdvertsRepo struct {
	*sqlite3.Sqlite
}

func NewAdvertsRepo(sq *sqlite3.Sqlite) *AdvertsRepo {
	return &AdvertsRepo{sq}
}

func (ar *AdvertsRepo) Store(ctx context.Context, adv *entity.Advert) error {
	tx, err := ar.DB.Begin()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Store - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	err = ar.storeAdvert(ctx, tx, adv)
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Store - %w", err)
	}

	for i := 0; i < len(adv.PhotosUrls); i++ {
		err := ar.storeUrl(ctx, tx, adv.Id, adv.PhotosUrls[i])
		if err != nil {
			return fmt.Errorf("AdvertsRepo - Store - %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Store - Commit: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) storeAdvert(ctx context.Context, tx *sql.Tx, adv *entity.Advert) error {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO adverts(name, description, price, photo_url, created_at) 
		values(?, ?, ?, ?, ?)`,
		adv.Name, adv.Description, adv.Price, adv.PhotosUrls[0], adv.CreatedAt)
	if err != nil {
		return fmt.Errorf("storeAdvert - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("storeAdvert - RowsAffected: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("storeAdvert - LastInsertId: %w", err)
	}
	adv.Id = id

	return nil
}

func (ar *AdvertsRepo) storeUrl(ctx context.Context, tx *sql.Tx, advId int64, url string) error {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO photo_urls(advert_id, url) values(?, ?)`,
		advId, url)
	if err != nil {
		return fmt.Errorf("storeUrl - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("storeUrl - RowsAffected: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) GetById(ctx context.Context, id int64) (entity.Advert, error) {
	advert := entity.Advert{}

	tx, err := ar.DB.Begin()
	if err != nil {
		return advert, fmt.Errorf("AdvertsRepo - GetById - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	row := tx.QueryRowContext(ctx,
		`SELECT 
            id, name, description, price, photo_url, created_at
        FROM adverts             
        WHERE id = ?`, id)

	var description sql.NullString
	var price sql.NullInt64
	var url sql.NullString

	err = row.Scan(&advert.Id, &advert.Name, &description, &price, &url, &advert.CreatedAt)
	if err != nil {
		return advert, fmt.Errorf("AdvertsRepo - GetById - Scan: %w", err)
	}

	advert.Description = description.String
	advert.Price = price.Int64
	advert.MainPhotoUrl = url.String

	urls, err := ar.getUrls(ctx, tx, advert.Id)
	if err != nil {
		return advert, fmt.Errorf("AdvertsRepo - GetById - %w", err)
	}

	advert.PhotosUrls = append(advert.PhotosUrls, urls...)

	err = tx.Commit()
	if err != nil {
		return advert, fmt.Errorf("AdvertsRepo - GetById - Commit: %w", err)
	}

	return advert, nil

}

func (ar *AdvertsRepo) Fetch(ctx context.Context) ([]entity.Advert, error) {
	adverts := []entity.Advert{}

	rows, err := ar.DB.QueryContext(ctx,
		`SELECT 
            name, price, photo_url, created_at
            FROM adverts             
        `)
	if err != nil {
		return adverts, fmt.Errorf("AdvertsRepo - Fetch - QueryContext: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var advert entity.Advert
		var price sql.NullInt64
		var url sql.NullString

		err = rows.Scan(&advert.Name, &price, &url, &advert.CreatedAt)
		if err != nil {
			return adverts, fmt.Errorf("AdvertsRepo - Fetch - Scan: %w", err)
		}

		advert.Price = price.Int64
		advert.MainPhotoUrl = url.String

		adverts = append(adverts, advert)
	}

	return adverts, nil
}

func (ar *AdvertsRepo) getUrls(ctx context.Context, tx *sql.Tx,
	advId int64) ([]string, error) {
	urls := []string{}
	rows, err := tx.QueryContext(ctx,
		`SELECT url
                FROM photo_urls
                WHERE advert_id = ?              
        `, advId)
	if err != nil {
		return urls, fmt.Errorf("getUrls - Exec: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var url sql.NullString

		err = rows.Scan(&url)
		if err != nil {
			return urls, fmt.Errorf("getUrls - Scan: %w", err)
		}
		urls = append(urls, url.String)
	}

	return urls, nil
}

func (ar *AdvertsRepo) Update(ctx context.Context, adv entity.Advert) error {
	tx, err := ar.DB.Begin()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx,
		`UPDATE adverts 
                SET name = ?, description = ?, price = ?, photo_url = ?
                WHERE id = ? 
                `, adv.Name, adv.Description, adv.Price, adv.MainPhotoUrl, adv.Id)

	if err != nil {
		return fmt.Errorf("AdvertsRepo - Update - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return entity.ErrItemNotExists
	}

	err = ar.updateUrls(ctx, tx, adv)

	if err != nil {
		return fmt.Errorf("AdvertsRepo - Update - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - Commit: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) updateUrls(ctx context.Context, tx *sql.Tx,
	adv entity.Advert) error {
	err := ar.deleteUrls(ctx, tx, adv.Id)
	if err != nil {
		return fmt.Errorf("updateUrls - %w", err)
	}
	if len(adv.PhotosUrls) != 0 {
		for i := 0; i < len(adv.PhotosUrls); i++ {
			err := ar.storeUrl(ctx, tx, adv.Id, adv.PhotosUrls[i])
			if err != nil {
				return fmt.Errorf("updateUrls - %w", err)
			}
		}
	}

	return nil
}

func (ar *AdvertsRepo) Delete(ctx context.Context, id int64) error {
	tx, err := ar.DB.Begin()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	err = ar.deleteAdvert(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - %w", err)
	}

	err = ar.deleteUrls(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - Commit: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) deleteAdvert(ctx context.Context, tx *sql.Tx, id int64) error {
	res, err := tx.ExecContext(ctx,
		`DELETE FROM adverts
                WHERE id = ?
                `, id)

	if err != nil {
		return fmt.Errorf("deleteAdvert - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return entity.ErrItemNotExists
	}

	return nil
}

func (ar *AdvertsRepo) deleteUrls(ctx context.Context, tx *sql.Tx, id int64) error {
	_, err := tx.ExecContext(ctx,
		`DELETE FROM photo_urls
                WHERE advert_id = ?
                `, id)

	if err != nil {
		return fmt.Errorf("deleteUrls - ExecContext: %w", err)
	}

	return nil
}
