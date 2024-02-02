package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nem0z/wiki-pathfinder/crawler"
	"github.com/nem0z/wiki-pathfinder/storage"
)

func generateRandomString(random *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}

	return string(result)
}

func GenerateRandomArticles(X int) []*storage.Article {
	randSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randSource)

	var articles []*storage.Article

	for i := 0; i < X; i++ {
		article := &storage.Article{
			Link:  generateRandomString(random, 10),
			Title: generateRandomString(random, 20),
		}

		articles = append(articles, article)
	}

	return articles
}

func worker(db *storage.DB, chArticle chan *storage.Article, chEnd chan bool) {
	for article := range chArticle {
		toArticles := GenerateRandomArticles(2500)
		db.CreateLinks(article.Id, toArticles...)
		chEnd <- true
	}
}

func main2() {
	db, _ := storage.Init("local.db")

	dbOut, _ := storage.Init("test.db")

	queue := crawler.InitQueue(db)
	ch := make(chan bool)
	chArticle := make(chan *storage.Article)

	for i := 0; i < 1; i++ {
		go worker(dbOut, chArticle, ch)
	}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		article := queue.Consume()
		chArticle <- article
	}

	for i := 0; i < 1000; i++ {
		<-ch
	}

	duration := time.Since(start)
	fmt.Printf("Execution time: %s\n", duration)
}
