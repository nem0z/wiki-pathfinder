package main

import (
	"log"
	"os"
	"strconv"

	"github.com/nem0z/wiki-pathfinder/crawler"
	"github.com/nem0z/wiki-pathfinder/storage"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := os.Args
	nbCrawlers := 500
	if len(args) > 1 {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			nbCrawlers = n
		}
	}

	db, err := storage.Init("local.db")
	Handle(err)

	queue := crawler.InitQueue(db)
	ch := make(chan *crawler.CrawlerResp)

	for i := 0; i < nbCrawlers; i++ {
		crawler := crawler.New(queue, ch)
		go crawler.Work()
	}
	log.Printf("Started %v crawlers", nbCrawlers)

	for {
		select {
		case resp := <-ch:
			db.CreateLinks(resp.ParentId, resp.Childs...)
		}
	}
}
