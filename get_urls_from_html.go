package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlbody, rawBaseURL string) ([]string, error) {
	htmlReader := strings.NewReader(htmlbody)
	nodeTree, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}

	rawLinks := findLinks(nodeTree)

	if len(rawLinks) == 0 {
		return []string{}, err
	}

	baseUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	var resolvedLinks []string
	seen := make(map[string]bool)
	for _, rawLink := range rawLinks {
		if rawLink == "" {
			continue
		}

		parsedLink, err := url.Parse(rawLink)
		if err != nil {
			continue
		}

		resolvedLink := baseUrl.ResolveReference(parsedLink).String()
		if !seen[resolvedLink] {
			seen[resolvedLink] = true
			resolvedLinks = append(resolvedLinks, resolvedLink)
		}
	}

	return resolvedLinks, nil
}

func findLinks(node *html.Node) []string {
	var links []string
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, findLinks(c)...)
	}
	return links
}
