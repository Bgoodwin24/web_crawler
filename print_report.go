package main

import (
	"fmt"
	"sort"
)

type Page struct {
	URL   string
	Count int
}

func sortReport(pages map[string]int) []Page {
	var sortedPages []Page

	for url, count := range pages {
		sortedPages = append(sortedPages, Page{
			URL:   url,
			Count: count,
		})
	}

	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].Count == sortedPages[j].Count {
			return sortedPages[i].URL < sortedPages[j].URL
		}
		return sortedPages[i].Count > sortedPages[j].Count
	})

	return sortedPages
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf("=============================\n  REPORT for %s\n=============================\n", baseURL)

	sorted := sortReport(pages)

	for _, page := range sorted {
		fmt.Printf("Found %d internal links to %s\n", page.Count, page.URL)
	}
}
