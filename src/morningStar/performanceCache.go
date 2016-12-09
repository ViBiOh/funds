package morningStar

type cacheRequest struct {
	key     string
	entries chan *performance
}

func cacheServer(chan<- *cacheRequest) {
	cache := make(map[string]*performance)

	for req := range cacheRequests {
		if req.key == `list` {
			for _, perf := range cache {
				req.entries <- perf
			}
			close(req.entry)
		} else if req.key != `` {
			if entry, ok = cache[req.key]; ok {
				req.entries <- entry
			}
			close(req.entry)
		} else {
			for entry := range req.entries {
				cache[entry.ID] = entry
			}
		}
	}
}
