package main

import (
	"fmt"
	"time"

	"github.com/nem0z/wiki-pathfinder/crawler"
	"github.com/nem0z/wiki-pathfinder/storage"
)

func main() {
	db, _ := storage.Init("local.db")

	queue := crawler.InitQueue(db)
	ch := make(chan *crawler.CrawlerResp)

	start := time.Now()
	for i := 0; i < 50; i++ {
		crawler := crawler.New(queue, ch)
		go crawler.Work()
	}

	for i := 0; i < 10000; i++ {
		<-ch
	}

	duration := time.Since(start)
	fmt.Printf("Execution time: %s\n", duration)
}
