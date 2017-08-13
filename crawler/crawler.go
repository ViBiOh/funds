package crawler

import (
	"sync"
)

// MaxConcurrentFetcher count in parallel
const MaxConcurrentFetcher = 32

// Error contains ID in error and error desc
type Error struct {
	ID  []byte
	Err error
}

func concurrentCrawl(ids [][]byte, fetcher func([]byte) (interface{}, error), wg *sync.WaitGroup, results chan<- interface{}, errors chan<- *Error) {
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
				errors <- &Error{ID, err}
			}
		}(ID)
	}
}

// Crawl retrieve given ids by calling fetcher func in parallel
func Crawl(ids [][]byte, fetcher func([]byte) (interface{}, error)) (<-chan interface{}, <-chan *Error) {
	var wgFetch sync.WaitGroup
	wgFetch.Add(len(ids))

	results := make(chan interface{}, MaxConcurrentFetcher)
	errors := make(chan *Error, MaxConcurrentFetcher)

	go concurrentCrawl(ids, fetcher, &wgFetch, results, errors)

	go func() {
		wgFetch.Wait()
		close(results)
		close(errors)
	}()

	return results, errors
}
