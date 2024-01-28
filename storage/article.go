package storage

import (
	"database/sql"
	"errors"
)

type Article struct {
	Id    int64
	Link  string
	Title string
}

func NewArticle(link string, title string) *Article {
	return &Article{Id: 0, Link: link, Title: title}
}

func (db *DB) GetArticle(article *Article) (id int64, err error) {
	err = db.QueryRow("SELECT id FROM articles WHERE link = ?", article.Link).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

func (db *DB) CreateArticle(article *Article) (int64, error) {
	req := "INSERT INTO articles (link, title, processed) VALUES (?, ?, ?)"
	res, err := db.Exec(req, article.Link, article.Title, false)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (db *DB) GetOrCreateArticle(article *Article) (int64, error) {
	id, err := db.GetArticle(article)
	if err != nil {
		return 0, err
	}

	if id != 0 {
		return id, nil
	}

	return db.CreateArticle(article)
}

func (db *DB) SetArticleProcessed(id int64) error {
	_, err := db.Exec("UPDATE articles SET processed = 1 WHERE id = ?", id)
	return err
}

func loadArticle(rows *sql.Rows) (*Article, error) {
	article := &Article{}
	err := rows.Scan(&article.Id, &article.Title, &article.Link)
	return article, err
}

func loadArticles(rows *sql.Rows) (articles []*Article, err error) {
	for rows.Next() {
		article, err := loadArticle(rows)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, err
}

func (db *DB) LoadArticles(count int, processed bool) ([]*Article, error) {
	query := `
		SELECT id, title, link
		FROM articles
		WHERE processed = ?
		LIMIT ?;
	`

	rows, err := db.Query(query, processed, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return loadArticles(rows)
}
