package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const (
	createTableArticles = `
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY,
			title VARCHAR(255),
			link VARCHAR(255) UNIQUE,
			processed bool
		)
	`

	createTableLinks = `
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY,
			parent INTEGER,
			child INTEGER,
			UNIQUE (parent, child),
			FOREIGN KEY (parent) REFERENCES articles(id) ON DELETE CASCADE,
			FOREIGN KEY (child) REFERENCES articles(id) ON DELETE CASCADE
		);
	`

	insertOriginArticle = `
		INSERT INTO articles ("title", "link", "processed")
		SELECT 'Strasbourg', '/wiki/Strasbourg', '0'
		WHERE NOT EXISTS (SELECT 1 FROM articles LIMIT 1);
	`
)

type DB struct {
	*sql.DB
}

type Tx struct {
	*sql.Tx
}

func Init(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableArticles)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableLinks)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(insertOriginArticle)
	return &DB{db}, err
}

func (db *DB) BeginTx() (*Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}
