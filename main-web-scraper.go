package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func scrape(url string) ([]Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var items []Item
	doc.Find("a.storylink").Each(func(_ int, s *goquery.Selection) {
		title := s.Text()
		href, _ := s.Attr("href")
		items = append(items, Item{Title: title, URL: href})
	})
	return items, nil
}

func main() {
	singleURL := flag.String("url", "https://news.ycombinator.com", "Single page URL to scrape")
	urlsFlag := flag.String("urls", "", "Comma-separated list of URLs to scrape")
	out := flag.String("out", "results.json", "Output JSON file")
	flag.Parse()

	var urls []string
	if *urlsFlag != "" {
		urls = strings.Split(*urlsFlag, ",")
	} else {
		urls = []string{*singleURL}
	}

	resultsChan := make(chan []Item, len(urls))
	var wg sync.WaitGroup

	// For each URL, start a goroutine and add the task to the WaitGroup.
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			items, err := scrape(url)
			if err != nil {
				log.Printf("Error scraping %s: %v\n", url, err)
				resultsChan <- []Item{}
				return
			}
			resultsChan <- items
		}(u)
	}

	wg.Wait()
	close(resultsChan)

	var allItems []Item
	for part := range resultsChan {
		allItems = append(allItems, part...)
	}

	data, err := json.MarshalIndent(allItems, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}
	if err := os.WriteFile(*out, data, 0644); err != nil {
		log.Fatalf("Write file failed: %v", err)
	}

	fmt.Printf("Scraped %d items from %d URL(s). Saved to %s\n", len(allItems), len(urls), *out)
}
