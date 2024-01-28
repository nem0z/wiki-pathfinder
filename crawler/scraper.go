package crawler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/nem0z/wiki-pathfinder/storage"
)

type Scraper struct {
	*colly.Collector
}

func NewScraper() *Scraper {
	return &Scraper{colly.NewCollector()}
}

func isValidLink(link string) bool {
	return strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":")
}

func (s *Scraper) GetArticles(link string) (articles []*storage.Article, finalError error) {
	s.OnHTML("#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")

		if isValidLink(link) {
			articles = append(articles, storage.NewArticle(link, title))
		}
	})

	s.OnError(func(r *colly.Response, err error) {
		if err != nil {
			formatedError := fmt.Sprintf("Request URL : %v failed with response: %v\nError : %v", r.Request.URL, r, err)
			finalError = errors.New(formatedError)
		}
	})

	url := fmt.Sprintf("https://fr.wikipedia.org%v", link)
	err := s.Visit(url)
	if finalError != nil {
		finalError = err
	}

	return articles, finalError
}
