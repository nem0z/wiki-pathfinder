package main

import (
	"log"

	"github.com/nem0z/wiki-pathfinder/crawler"
	"github.com/nem0z/wiki-pathfinder/storage"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := storage.Init("local.db")
	Handle(err)

	queue := crawler.InitQueue(db)
	ch := make(chan *crawler.CrawlerResp)

	for i := 0; i < 10; i++ {
		crawler := crawler.New(queue, ch)
		go crawler.Work()
	}

	for {
		select {
		case resp := <-ch:
			db.CreateLinks(resp.ParentId, resp.Childs...)
		}
	}
}
