package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	createTableArticles = `
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255),
			link VARCHAR(255) UNIQUE,
			processed bool
		)
	`

	createTableLinks = `
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER AUTO_INCREMENT PRIMARY KEY,
			parent INTEGER,
			child INTEGER,
			UNIQUE (parent, child),
			FOREIGN KEY (parent) REFERENCES articles(id) ON DELETE CASCADE,
			FOREIGN KEY (child) REFERENCES articles(id) ON DELETE CASCADE
		);
	`

	insertOriginArticle = `
		INSERT INTO articles (title, link, processed)
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
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Format MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%v", user, password, host, port, name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check if the connection to the database is successful
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create the database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		return nil, err
	}

	// Switch to the specified database
	// _, err = db.Exec("USE " + name)
	// if err != nil {
	// 	return nil, err
	// }

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
