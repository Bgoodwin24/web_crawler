package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("not enough arguments, usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	if len(args) > 4 {
		fmt.Println("too many arguments, usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	baseURL := args[1]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	maxConcurrencyInt, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("maxConcurrency must be an integer")
		os.Exit(1)
	}

	maxPagesInt, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("maxPages must be an integer")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("error parsing URL: %v\n", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            parsedURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrencyInt),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPagesInt,
	}

	cfg.crawlPage(baseURL)

	cfg.wg.Wait()

	fmt.Println("Results:")
	for url, count := range cfg.pages {
		fmt.Printf("Found %s %d times\n", url, count)
	}
}
