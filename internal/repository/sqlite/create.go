package sqlite

import (
	"github.com/mrsubudei/adv-store-service/pkg/sqlite3"
)

func CreateDB(s *sqlite3.Sqlite) error {

	adverts := `
	CREATE TABLE IF NOT EXISTS adverts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		price INTEGER,
		photo_url TEXT
		);
	`

	_, err := s.DB.Exec(adverts)
	if err != nil {
		return err
	}

	urls := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		photo_url TEXT
		);
	`
	_, err = s.DB.Exec(urls)
	if err != nil {
		return err
	}
	urlReference := `
	CREATE TABLE IF NOT EXISTS url_reference (
		advert_id INTEGER,
		url_id INTEGER,
		PRIMARY KEY (advert_id, url_id),
		FOREIGN KEY (advert_id) REFERENCES adverts(id)
		FOREIGN KEY (url_id) REFERENCES urls(id)
		);
	`
	_, err = s.DB.Exec(urlReference)
	if err != nil {
		return err
	}

	return nil
}
