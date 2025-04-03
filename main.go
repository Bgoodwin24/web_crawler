package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func main() {
	args := os.Args
	if len(args)-1 < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args)-1 > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[1]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("error parsing URL: %v\n", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            parsedURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 10),
		wg:                 &sync.WaitGroup{},
	}

	cfg.crawlPage(baseURL)

	cfg.wg.Wait()

	fmt.Println("Results:")
	for url, count := range cfg.pages {
		fmt.Printf("Found %s %d times\n", url, count)
	}
}
