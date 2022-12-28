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
		referenceId, err := ar.storeUrls(ctx, tx, adv.PhotosUrls[i])
		if err != nil {
			return fmt.Errorf("AdvertsRepo - Store - %w", err)
		}
		err = ar.storeReferenceUrls(ctx, tx, adv.Id, referenceId)
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
		`INSERT INTO adverts(name, description, price, photo_url) values(?, ?, ?, ?)`,
		adv.Name, adv.Description, adv.Price, adv.MainPhotoUrl)
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

func (ar *AdvertsRepo) storeUrls(ctx context.Context, tx *sql.Tx, url string) (int64, error) {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO urls(photo_url) values(?)`,
		url)
	if err != nil {
		return 0, fmt.Errorf("storeUrls - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return 0, fmt.Errorf("storeUrls - RowsAffected: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("storeUrls - LastInsertId: %w", err)
	}

	return id, nil
}

func (ar *AdvertsRepo) storeReferenceUrls(ctx context.Context, tx *sql.Tx, advId, urlId int64) error {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO url_reference(advert_id, url_id) values(?)`,
		advId, urlId)
	if err != nil {
		return fmt.Errorf("storeReferenceUrls - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("storeReferenceUrls - RowsAffected: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) GetById(ctx context.Context, id int64) (entity.Advert, error) {
    
    tx, err := ar.DB.Begin()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - GetById - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

    advert := entity.Advert{}
    
	row := tx.QueryRowContext(ctx,
		`SELECT 
			id, name, description, price, photo_url
		FROM adverts		 
		WHERE id = ?`)
	
	var description sql.NullString
	var price sql.NullInt64
	var url sql.NullString

	err := row.Scan(ctx, &advert.Id, &advert.Name, &description, &price, &url)
	if err != nil {
		return advert, fmt.Errorf("AdvertsRepo - GetById - Scan: %w", err)
	}

	advert.Description = description.String
	advert.Price = price.Int64
	advert.MainPhotoUrl = url.String
	
	urls, err := ar.getUrls(ctx, tx, advert.Id)
	if err != nil {
	    return fmt.Errorf("AdvertsRepo - GetById - %w", err)
	}
    
    advert.PhotosUrls = append(PhotosUrls, urls...)
    
    err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - GetById - Commit: %w", err)
	}
	
	return advert, nil
	
	
}

func (ar *AdvertsRepo) Fetch(ctx context.Context) ([]entity.Advert, error) {
	adverts := []entity.Advert{}
	rows, err := ar.DB.QueryContext(ctx,
		`SELECT 
			id, name, description, price, photo_url
		FROM adverts		 
	`)
	if err != nil {
		return adverts, fmt.Errorf("AdvertsRepo - Fetch - QueryContext: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var advert entity.Advert
		var description sql.NullString
		var price sql.NullInt64
		var url sql.NullString

		err = rows.Scan(ctx, &advert.Id, &advert.Name, &description, &price, &url)
		if err != nil {
			return adverts, fmt.Errorf("AdvertsRepo - Fetch - Scan: %w", err)
		}

		advert.Description = description.String
		advert.Price = price.Int64
		advert.MainPhotoUrl = url.String

    	urls, err := ar.getUrls(ctx, tx, advert.Id)
	    if err != nil {
	        return fmt.Errorf("AdvertsRepo - Fetch - %w", err)
    	}
    
         advert.PhotosUrls = append(PhotosUrls, urls...)
    
		adverts = append(adverts, advert)
	}

	return adverts, nil
}

func (ar *AdvertsRepo) getUrlIds(ctx context.Context, tx *sql.Tx, advId int64) ([]int64, error) {
	urlIds := []int64{}
	rows, err := tx.QueryContext(ctx,
		`SELECT url_id
		FROM url_reference
		WHERE advert_id = ?		 
	`, advId)
	if err != nil {
		return urlIds, fmt.Errorf("AdvertsRepo - getUrlIds - Exec: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64

		err = rows.Scan(ctx, &id)
		if err != nil {
			return urlIds, fmt.Errorf("AdvertsRepo - getUrlIds - Scan: %w", err)
		}
		urlIds = append(urlIds, id)
	}

	return urlIds, nil
}

func (ar *AdvertsRepo) getUrls(ctx context.Context, tx *sql.Tx, advId int64) ([]string, error) {
	urls := []string{}
	rows, err := tx.QueryContext(ctx,
		`SELECT photo_url
		FROM urls
		WHERE id = (SELECT url_id FROM url_reference WHERE url_reference.advert_id = ?)		 
	`)
	if err != nil {
		return urls, fmt.Errorf("AdvertsRepo - getUrls - QueryContext: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var url sql.NullString

		err = rows.Scan(ctx, &url)
		if err != nil {
			return urls, fmt.Errorf("AdvertsRepo - getUrls - Scan: %w", err)
		}
		urls = append(urls, url.String)
	}

	return urls, nil
}

func (ar *AdvertsRepo) Update(ctx context.Context, adv entity.Advert) error {
	res, err := ar.DB.ExecContext(ctx,
		`UPDATE adverts 
		SET name = ?, description = ?, price = ?, photo_url = ?
		WHERE id = ? 
		`, adv.Name, adv.Description, adv.Price, adv.MainPhotoUrl, adv.Id)

	if err != nil {
		return fmt.Errorf("AdvertsRepo - Update - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("AdvertsRepo - Update - RowsAffected: %w", err)
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
		err = ar.deleteReferenceUrls(ctx, tx, id)
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
		return fmt.Errorf("deleteAdvert - RowsAffected: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) deleteUrls(ctx context.Context, tx *sql.Tx, id int64) error {
	res, err := tx.ExecContext(ctx,
		`DELETE FROM urls
		WHERE id = (SELECT url_id FROM url_reference WHERE url_reference.advert_id = ?)
		`, id)

	if err != nil {
		return fmt.Errorf("deleteUrls - ExecContext: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) deleteReferenceUrls(ctx context.Context, tx *sql.Tx, id int64) error {
	res, err := tx.ExecContext(ctx,
		`DELETE FROM url_reference
		WHERE advert_id = ?
		`, id)

	if err != nil {
		return fmt.Errorf("deleteReferenceUrls - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("deleteReferenceUrls - RowsAffected: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) DeleteUrl(ctx context.Context, url string) error {
	tx, err := ar.DB.Begin()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - DeleteUrl - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx,
		`DELETE FROM url_reference
		WHERE url_id = (SELECT id FROM urls WHERE urls.id = ?)
		`, url)

	if err != nil {
		return fmt.Errorf("deleteUrls - ExecContext #1: %w", err)
	}
	
	res, err := tx.ExecContext(ctx,
		`DELETE FROM urls
		WHERE id = ?
		`, url)

	if err != nil {
		return fmt.Errorf("deleteUrls - ExecContext #2: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - DeleteUrl - Commit: %w", err)
	}

	return nil
}
