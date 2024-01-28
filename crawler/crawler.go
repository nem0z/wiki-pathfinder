package crawler

import (
	"fmt"
	"time"

	"github.com/nem0z/wiki-pathfinder/storage"
)

type CrawlerResp struct {
	ParentId int64
	Childs   []*storage.Article
}

type Crawler struct {
	queue   *Queue
	scraper *Scraper
	ch      chan *CrawlerResp
}

func New(queue *Queue, ch chan *CrawlerResp) *Crawler {
	return &Crawler{queue, NewScraper(), ch}
}

func (c *Crawler) Work() {
	for {
		item := c.queue.Consume()
		if item == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Println("Item to process :", item.Title)

		articles, err := c.scraper.GetArticles(item.Link)
		if err != nil {
			fmt.Println("crawler.GetArticles :", err)
		} else {
			c.ch <- &CrawlerResp{item.Id, articles}
		}
	}
}
