package main

import (
	"fmt"
	"time"

	"github.com/nem0z/wiki-pathfinder/crawler"
	"github.com/nem0z/wiki-pathfinder/storage"
)

func main1() {
	db, _ := storage.Init("local.db")

	queue := crawler.InitQueue(db)
	queue.SetMax(1000000)
	queue.Refill()

	start := time.Now()
	for i := crawler.MinQueueElements; i < crawler.MaxQueueElements; i++ {
		queue.Consume()
	}
	duration := time.Since(start)
	fmt.Printf("Execution time: %s\n", duration)
}
