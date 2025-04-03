package main

import (
	"fmt"
	"log"
	"strings"
)

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	count, exists := cfg.pages[normalizedURL]
	if exists {
		cfg.pages[normalizedURL] = count + 1
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.wg.Add(1)
	cfg.concurrencyControl <- struct{}{}
	defer cfg.wg.Done()
	defer func() { <-cfg.concurrencyControl }()

	if !strings.Contains(rawCurrentURL, cfg.baseURL.String()) {
		return
	}

	normalCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("error normalizing URL: %v", err)
		return
	}

	isFirstVisit := cfg.addPageVisit(normalCurrent)
	if !isFirstVisit {
		return
	}

	fmt.Printf("Currently crawling: %s\n", rawCurrentURL)
	currentHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML from %s: %v", rawCurrentURL, err)
		return
	}

	links, err := getURLsFromHTML(currentHTML, rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting links: %s from: %v", links, err)
	}

	for _, link := range links {
		go cfg.crawlPage(link)
	}
}
