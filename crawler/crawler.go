package crawler

import (
	"log"
	"sync"
)

// MaxConcurrentFetcher count in parallel
const MaxConcurrentFetcher = 32

func concurrentCrawl(ids [][]byte, fetcher func([]byte) (interface{}, error), wg *sync.WaitGroup, results chan<- interface{}, errors chan<- []byte) {
	tokens := make(chan int, MaxConcurrentFetcher)

	clearSemaphores := func() {
		wg.Done()
		<-tokens
	}

	for _, ID := range ids {
		tokens <- 1

		go func(ID []byte) {
			defer clearSemaphores()

			if result, err := fetcher(ID); err == nil {
				results <- result
			} else {
				log.Print(err)
				errors <- ID
			}
		}(ID)
	}
}

// Crawl retrieve given ids by calling fetcher func in parallel
func Crawl(ids [][]byte, fetcher func([]byte) (interface{}, error)) (<-chan interface{}, <-chan []byte) {
	var wgFetch sync.WaitGroup
	wgFetch.Add(len(ids))

	results := make(chan interface{}, MaxConcurrentFetcher)
	errors := make(chan []byte, MaxConcurrentFetcher)

	go concurrentCrawl(ids, fetcher, &wgFetch, results, errors)

	go func() {
		wgFetch.Wait()
		close(results)
		close(errors)
	}()

	return results, errors
}
