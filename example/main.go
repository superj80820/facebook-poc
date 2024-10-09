package main

import (
	"encoding/json"
	"os"
	"time"

	crawler "github.com/superj80820/facebook-poc"
	"github.com/superj80820/facebook-poc/fetcher"
	"github.com/superj80820/facebook-poc/formatter"
	"github.com/superj80820/facebook-poc/parser"
)

func main() {
	// args
	url := "https://www.facebook.com/anuetw/"
	startDateStr := "2024-10-05"

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		panic(err)
	}

	// cerate repository
	postFetcher := fetcher.NewFetcher(url)
	formatter := formatter.NewFormatter()
	parser := parser.NewParser()

	// create use case
	crawler := crawler.NewCrawler(postFetcher, formatter, parser)

	// business logic
	posts, err := crawler.FetchPagePosts(startDate, time.Now())
	if err != nil {
		panic(err)
	}

	// write to file
	m, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("./" + time.Now().String() + ".json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(m)
	if err != nil {
		panic(err)
	}
}
