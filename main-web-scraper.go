package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	url := flag.String("url", "https://news.ycombinator.com", "Page URL to scrape")
	out := flag.String("out", "results.json", "Output file")
	flag.Parse()

	items, err := scrape(*url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, data, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved %d items to %s\n", len(items), *out)
}
