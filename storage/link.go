package storage

import (
	"database/sql"
	"errors"
	"log"
)

func (tx *Tx) GetLink(parentId int64, childId int64) (id int64, err error) {
	err = tx.QueryRow("SELECT id FROM links WHERE parent = ? AND child = ?", parentId, childId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

func (tx *Tx) CreateLink(parentId int64, childId int64) (int64, error) {
	req := "INSERT INTO links (parent, child) VALUES (?, ?)"
	res, err := tx.Exec(req, parentId, childId)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (tx *Tx) GetOrCreateLink(parentId int64, childId int64) (int64, error) {
	id, err := tx.GetLink(parentId, childId)
	if err != nil {
		return 0, err
	}

	if id != 0 {
		return id, nil
	}

	return tx.CreateLink(parentId, childId)
}

func (tx *Tx) CreateLinks(parentId int64, childs ...*Article) {
	for _, article := range childs {
		id, err := tx.GetOrCreateArticle(article)
		if err != nil {
			log.Printf("CreateLinks => GetOrCreateArticle (%v) : %v\n", article.Link, err)
			continue
		}

		_, err = tx.GetOrCreateLink(parentId, id)
		if err != nil {
			log.Printf("CreateLinks => GetOrCreateLink (%v) : %v\n", id, err)
		}
	}
}

func (db *DB) CreateLinks(parentId int64, childs ...*Article) {
	tx, err := db.BeginTx()
	if err != nil {
		return
	}

	defer func() {
		err = tx.Commit()
		if err != nil {
			err = tx.SetArticleProcessed(parentId)
			if err != nil {
				log.Printf("CreateLinks => GetOrCreateArticle (%v) : %v\n", parentId, err)
			}
		}
	}()
}
