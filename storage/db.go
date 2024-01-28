package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

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

func Init(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableArticles)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableLinks)

	_, err = db.Exec(insertOriginArticle)
	return &DB{db}, nil
}
