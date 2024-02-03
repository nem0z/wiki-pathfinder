package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nem0z/wiki-pathfinder/storage/sharding"
)

func main() {
	shardMap, err := sharding.Init("links", 100)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()

	for i := 0; i < 1000000; i++ {
		shard, err := shardMap.GetShard(i)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(shard)
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time: %s\n", elapsedTime)
}
