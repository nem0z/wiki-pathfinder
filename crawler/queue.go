package crawler

import (
	"fmt"
	"log"
	"sync"

	"github.com/nem0z/wiki-pathfinder/storage"
)

const (
	minQueueElements = 500
	maxQueueElements = 5000
)

type Queue struct {
	arr []*storage.Article
	mu  sync.Mutex
	min int
	max int
	db  *storage.DB
}

func InitQueue(db *storage.DB) *Queue {
	q := &Queue{
		arr: []*storage.Article{},
		mu:  sync.Mutex{},
		min: minQueueElements,
		max: maxQueueElements,
		db:  db,
	}

	q.Refill()
	return q
}

func (q *Queue) Consume() *storage.Article {
	q.mu.Lock()

	if len(q.arr) < q.min {
		defer q.Refill()
	}

	defer q.mu.Unlock()

	if len(q.arr) == 0 {
		return nil
	}

	item := q.arr[0]
	q.arr = q.arr[1:]

	return item
}

func (q *Queue) Refill() {
	articles, err := q.db.LoadArticles(q.max, false)
	if err != nil {
		log.Println("Queue.Refill :", err)
	}

	fmt.Printf("Refill with %v articles\n", len(articles))
	q.mu.Lock()
	defer q.mu.Unlock()
	q.arr = articles
}
