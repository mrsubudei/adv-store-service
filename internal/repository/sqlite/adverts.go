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
		return fmt.Errorf("AdvertsRepo - BeginStore - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	err = ar.storeAdvert(ctx, tx, adv)
	if err != nil {
		return fmt.Errorf("AdvertsRepo - BeginStore - %w", err)
	}

	for i := 0; i < len(adv.PhotosUrls); i++ {
		referenceId, err := ar.storeUrls(ctx, tx, adv.PhotosUrls[i])
		if err != nil {
			return fmt.Errorf("AdvertsRepo - BeginStore - %w", err)
		}
		err = ar.storeReferenceUrls(ctx, tx, adv.Id, referenceId)
		if err != nil {
			return fmt.Errorf("AdvertsRepo - BeginStore - %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - BeginStore - Commit: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) storeAdvert(ctx context.Context, tx *sql.Tx, adv *entity.Advert) error {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO adverts(name, description, price, photo_url) values(?, ?, ?, ?)`,
		adv.Name, adv.Description, adv.Price, adv.MainPhotoUrl)
	if err != nil {
		return fmt.Errorf("store - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("store - RowsAffected: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("store - LastInsertId: %w", err)
	}
	adv.Id = id

	return nil
}

func (ar *AdvertsRepo) storeUrls(ctx context.Context, tx *sql.Tx, url string) (int64, error) {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO urls(photo_url) values(?)`,
		url)
	if err != nil {
		return 0, fmt.Errorf("StoreUrls - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return 0, fmt.Errorf("StoreUrls - RowsAffected: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("StoreUrls - LastInsertId: %w", err)
	}

	return id, nil
}

func (ar *AdvertsRepo) storeReferenceUrls(ctx context.Context, tx *sql.Tx, advId, urlId int64) error {
	res, err := tx.ExecContext(ctx,
		`INSERT INTO url_reference(advert_id, url_id) values(?)`,
		advId, urlId)
	if err != nil {
		return fmt.Errorf("StoreUrlReferences - Exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("StoreUrlReferences - RowsAffected: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) GetById(ctx context.Context, id int64) (entity.Advert, error) {
	advert := entity.Advert{}
	row := ar.DB.QueryRowContext(ctx,
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
		return adverts, fmt.Errorf("AdvertsRepo - Fetch - Exec: %w", err)
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

		adverts = append(adverts, advert)
	}

	return adverts, nil
}

func (ar *AdvertsRepo) GetUrlIds(ctx context.Context, adv entity.Advert) ([]int64, error) {
	urlIds := []int64{}
	rows, err := ar.DB.QueryContext(ctx,
		`SELECT url_id
		FROM url_reference
		WHERE advert_id = ?		 
	`, adv.Id)
	if err != nil {
		return urlIds, fmt.Errorf("AdvertsRepo - GetUrlIds - Exec: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64

		err = rows.Scan(ctx, &id)
		if err != nil {
			return urlIds, fmt.Errorf("AdvertsRepo - GetUrlIds - Scan: %w", err)
		}
		urlIds = append(urlIds, id)
	}

	return urlIds, nil
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

func (ar *AdvertsRepo) Delete(ctx context.Context, adv entity.Advert) error {
	tx, err := ar.DB.Begin()
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - Begin: %w", err)
	}
	defer func() {
		err = tx.Rollback()
	}()

	err = ar.deleteAdvert(ctx, tx, adv.Id)
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - %w", err)
	}

	urlIds, err := ar.GetUrlIds(ctx, adv)
	if err != nil {
		return fmt.Errorf("AdvertsRepo - Delete - %w", err)
	}

	for i := 0; i < len(urlIds); i++ {
		err = ar.deleteUrls(ctx, tx, urlIds[i])
		if err != nil {
			return fmt.Errorf("AdvertsRepo - Delete - %w", err)
		}
		err = ar.deleteReferenceUrls(ctx, tx, urlIds[i])
		if err != nil {
			return fmt.Errorf("AdvertsRepo - Delete - %w", err)
		}
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
		WHERE id = ?
		`, id)

	if err != nil {
		return fmt.Errorf("deleteUrls - ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if affected != 1 || err != nil {
		return fmt.Errorf("deleteUrls - RowsAffected: %w", err)
	}

	return nil
}

func (ar *AdvertsRepo) deleteReferenceUrls(ctx context.Context, tx *sql.Tx, id int64) error {
	res, err := tx.ExecContext(ctx,
		`DELETE FROM url_reference
		WHERE id = ?
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
