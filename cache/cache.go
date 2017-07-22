package cache

// Content stored in cache
type Content interface {
	GetID() string
}

// Request define information for requesting cache
type Request struct {
	key     string
	entries chan Content
}

// Get gets entry from cache connected to channel
func Get(ch chan<- Request, key string) Content {
	req := Request{key: key, entries: make(chan Content)}
	ch <- req

	return <-req.entries
}

// Push pushes entry to cache connected to channel
func Push(ch chan<- Request, entries []Content) {
	req := Request{key: `push`, entries: make(chan Content, 0)}
	ch <- req

	for _, entry := range entries {
		req.entries <- entry
	}

	close(req.entries)
}

// List dumps to output channel entries from cache connected to channel
func List(ch chan<- Request) <-chan Content {
	results := make(chan Content, 0)
	ch <- Request{key: `list`, entries: results}

	return results
}

// Server start a cache and listen for request from given channel
func Server(ch <-chan Request, size int) {
	cache := make(map[string]Content, size)

	for req := range ch {
		if req.key == `list` {
			for _, perf := range cache {
				req.entries <- perf
			}
			close(req.entries)
		} else if req.key == `push` {
			for entry := range req.entries {
				cache[entry.GetID()] = entry
			}
		} else if req.key != `` {
			if entry, ok := cache[req.key]; ok {
				req.entries <- entry
			}
			close(req.entries)
		}
	}
}
