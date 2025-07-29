package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

type pageCount struct {
	url   string
	count int
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if count, exists := cfg.pages[normalizedURL]; exists {
		cfg.pages[normalizedURL] = count + 1
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}

func (cfg *config) shouldContinueCrawling() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.pages) < cfg.maxPages
}

func normalizeURL(uri string) (string, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	cleanPath := strings.TrimSuffix(parsedURL.Path, "/")
	cleanHost := strings.TrimSuffix(parsedURL.Host, "/")
	modifiedURL := cleanHost + cleanPath
	return modifiedURL, err
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}
	var urls []string
	traversNodes(doc, &urls, rawBaseURL)
	return urls, nil
}

func traversNodes(node *html.Node, urls *[]string, rawBaseURL string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				parsedURL, err := url.Parse(attr.Val)
				if err != nil {
					break
				}
				baseURL, err := url.Parse(rawBaseURL)
				if err != nil {
					break
				}
				absoluteURL := baseURL.ResolveReference(parsedURL)
				*urls = append(*urls, absoluteURL.String())
				break
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traversNodes(child, urls, rawBaseURL)
	}
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("неправильний content-type: %s, очікувався text/html", contentType)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(htmlBytes), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	if !cfg.shouldContinueCrawling() {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if cfg.baseURL.Host != currentURL.Host {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	isFirst := cfg.addPageVisit(normalizedURL)
	if !isFirst {
		return
	}

	fmt.Printf("crawling: %s\n", rawCurrentURL)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Помилка при отриманні HTML: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Помилка при пошуку URL: %v\n", err)
		return
	}

	for _, url := range urls {
		if !cfg.shouldContinueCrawling() {
			break
		}

		cfg.wg.Add(1)
		go func(urlCopy string) {
			defer cfg.wg.Done()
			defer func() { <-cfg.concurrencyControl }()

			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(urlCopy)
		}(url)
	}
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Printf("=============================\n")
	fmt.Printf("  REPORT for %s\n", baseURL)
	fmt.Printf("=============================\n")

	pageCountSlice := make([]pageCount, 0, len(pages))
	for url, count := range pages {
		pageCountSlice = append(pageCountSlice, pageCount{
			url:   url,
			count: count,
		})
	}

	sort.Slice(pageCountSlice, func(i, j int) bool {
		if pageCountSlice[i].count != pageCountSlice[j].count {
			return pageCountSlice[i].count > pageCountSlice[j].count

		}
		return pageCountSlice[i].url < pageCountSlice[j].url
	})

	for _, pc := range pageCountSlice {
		fmt.Printf("Found %d internal links to https://%s\n", pc.count, pc.url)
	}
}
