package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	command := os.Args
	if len(command) != 4 {
		fmt.Println("Usage: crawler URL maxConcurrency maxPages")
		fmt.Println("Example: crawler https://example.com 3 10")
		os.Exit(1)
	}
	if len(command) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	maxConcurrency, err := strconv.Atoi(command[2])
	if err != nil || maxConcurrency < 1 {
		fmt.Println("maxConcurrency має бути позитивним числом")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(command[3])
	if err != nil || maxPages < 1 {
		fmt.Println("maxPages має бути позитивним числом")
		os.Exit(1)
	}
	baseURL, err := url.Parse(command[1])
	if err != nil {
		fmt.Printf("Помилка парсингу URL: %v\n", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
	fmt.Println("starting crawl of:", command[1])

	cfg.wg.Add(1)
	go func() {
		defer cfg.wg.Done()
		defer func() { <-cfg.concurrencyControl }()
		cfg.concurrencyControl <- struct{}{}
		cfg.crawlPage(command[1])
	}()

	cfg.wg.Wait()

	printReport(cfg.pages, command[1])
}
