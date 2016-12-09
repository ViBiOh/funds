package morningStar

type cacheRequest struct {
	key   string
	size  int
	entry chan *performance
}

func cacheServer(chan<- *cacheRequest) {
	cache := make(map[string]*performance)

	for req := range cacheRequests {
		if req.key == `list` {
			req.size = len(cache)
			for _, perf := range cache {
				req.entry <- perf
			}
			close(req.entry)
		} else if req.key != `` {
			if entry, ok = cache[req.key]; ok {
				req.entry <- entry
			}
			close(req.entry)
		} else {
			for entry := range req.entry {
				cache[entry.ID] = entry
			}
		}
	}
}
