package morningStar

type cacheRequest struct {
	key     string
	entries chan *performance
}

func getCache(ch chan<- *cacheRequest, key string) *performance {
	req := cacheRequest{key: key, entries: make(chan *performance)}
	ch <- &req

	return <-req.entries
}

func pushCache(ch chan<- *cacheRequest, entries []*performance) {
	req := cacheRequest{key: `push`, entries: make(chan *performance, 0)}
	ch <- &req

	for _, entry := range entries {
		req.entries <- entry
	}

	close(req.entries)
}

func listCache(ch chan<- *cacheRequest) <-chan *performance {
	results := make(chan *performance, 0)
	ch <- &cacheRequest{key: `list`, entries: results}

	return results
}

func cacheServer(ch <-chan *cacheRequest, size int) {
	cache := make(map[string]*performance, size)

	for req := range ch {
		if req.key == `list` {
			for _, perf := range cache {
				req.entries <- perf
			}
			close(req.entries)
		} else if req.key == `push` {
			for entry := range req.entries {
				cache[entry.ID] = entry
			}
		} else if req.key != `` {
			if entry, ok := cache[req.key]; ok {
				req.entries <- entry
			}
			close(req.entries)
		}
	}
}
