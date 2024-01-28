package storage

import (
	"database/sql"
	"errors"
	"log"
)

func (db *DB) GetLink(parentId int64, childId int64) (id int64, err error) {
	err = db.QueryRow("SELECT id FROM links WHERE parent = ? AND child = ?", parentId, childId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

func (db *DB) CreateLink(parentId int64, childId int64) (int64, error) {
	req := "INSERT INTO links (parent, child) VALUES (?, ?)"
	res, err := db.Exec(req, parentId, childId)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (db *DB) GetOrCreateLink(parentId int64, childId int64) (int64, error) {
	id, err := db.GetLink(parentId, childId)
	if err != nil {
		return 0, err
	}

	if id != 0 {
		return id, nil
	}

	return db.CreateLink(parentId, childId)
}

func (db *DB) CreateLinks(parentId int64, childs ...*Article) {
	defer db.SetArticleProcessed(parentId)

	for _, article := range childs {
		id, err := db.GetOrCreateArticle(article)
		if err != nil {
			log.Printf("CreateLinks => GetOrCreateArticle (%v) : %v\n", article.Link, err)
			continue
		}

		_, err = db.GetOrCreateLink(parentId, id)
		if err != nil {
			log.Printf("CreateLinks => GetOrCreateLink (%v) : %v\n", id, err)
		}
	}
}
