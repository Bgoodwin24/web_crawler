package main

import (
	"fmt"
	"log"
	"strings"
)

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, exists := cfg.pages[normalizedURL]
	if exists {
		return false
	}

	cfg.pages[normalizedURL] = 0
	return true
}

func (cfg *config) incrementLinkCount(normalizedURL string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.pages[normalizedURL]++
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	cfg.concurrencyControl <- struct{}{}
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
		if strings.Contains(link, cfg.baseURL.String()) {
			normalizedLink, err := normalizeURL(link)
			if err != nil {
				continue
			}

			cfg.incrementLinkCount(normalizedLink)

			cfg.mu.Lock()
			shouldCrawl := len(cfg.pages) < cfg.maxPages
			cfg.mu.Unlock()

			if shouldCrawl {
				cfg.wg.Add(1)
				go cfg.crawlPage(link)
			}
		}
	}
}
