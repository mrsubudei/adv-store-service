package sqlite

import (
	"fmt"

	"github.com/mrsubudei/adv-store-service/pkg/sqlite3"
)

func CreateDB(s *sqlite3.Sqlite) error {

	adverts := `
	CREATE TABLE IF NOT EXISTS adverts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		price INTEGER,
		photo_url TEXT,
		created_at TEXT
		);
	`

	_, err := s.DB.Exec(adverts)
	if err != nil {
		return fmt.Errorf("CreateDB - %w", err)
	}

	urls := `
	CREATE TABLE IF NOT EXISTS photo_urls (
		advert_id INTEGER,
		url TEXT,
		PRIMARY KEY (advert_id, url),
		FOREIGN KEY (advert_id) REFERENCES adverts(id)
		);
	`
	_, err = s.DB.Exec(urls)
	if err != nil {
		return fmt.Errorf("CreateDB - %w", err)
	}

	return nil
}
