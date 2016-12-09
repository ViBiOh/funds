package morningStar

const bufferSize = 8

type cacheRequest struct {
	key     string
	entries chan *performance
}

func getCache(ch chan<- *cacheRequest, key string) *performance {
	req := cacheRequest{key: key, entries: make(chan *performance)}
	ch <- &req
	
	return <- req.entries
}

func pushCache(ch chan<- *cacheRequest, entry *performance) {
	req := cacheRequest{entries: make(chan *performance)}
	ch <- &req
	req.entries <- entry

	close(req.entries)
}

func loadCache(ch chan<- *cacheRequest, entries []*performance) {
	req := cacheRequest{entries: make(chan *performance, bufferSize)}
	ch <- &req
	
	for _, entry := range entries {
		req.entries <- entry
	}

	close(req.entries)
}

func listCache(ch chan<- *cacheRequest) <-chan *performance {
	results := make(chan *performance, bufferSize)
	ch <- &cacheRequest{key: `list`, entries: results}
	
	return results
}

func cacheServer(chan<- *cacheRequest) {
	cache := make(map[string]*performance)

	for req := range cacheRequests {
		if req.key == `list` {
			for _, perf := range cache {
				req.entries <- perf
			}
			close(req.entries)
		} else if req.key != `` {
			if entry, ok := cache[req.key]; ok {
				req.entries <- entry
			}
			close(req.entries)
		} else {
			for entry := range req.entries {
				cache[entry.ID] = entry
			}
		}
	}
}
