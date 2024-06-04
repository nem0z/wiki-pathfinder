package crawler

import (
	"log"
	"time"

	"github.com/nem0z/wiki-pathfinder/storage"
)

type CrawlerResp struct {
	ParentId int64
	Childs   []*storage.Article
}

type Crawler struct {
	queue *Queue
	ch    chan *CrawlerResp
}

func New(queue *Queue, ch chan *CrawlerResp) *Crawler {
	return &Crawler{queue, ch}
}

func (c *Crawler) Work() {
	for {
		item := c.queue.Consume()
		if item == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Println("Item to process :", item.Title)
		scraper := NewScraper()
		articles, err := scraper.GetArticles(item.Link)
		if err != nil {
			log.Println("crawler.GetArticles :", err)
		} else {
			c.ch <- &CrawlerResp{item.Id, articles}
		}
	}
}
