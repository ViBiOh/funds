package morningStar

type cacheRequest struct {
	key   string
	value *performance
	list  []*performance
	ready chan int
}

func performanceCacheServer(chan<- *cacheRequest) {
	cache := make(map[string]*performance)

	var ready = func(request *cacheRequest) {
		if request.ready != nil {
			close(request.ready)
		}
	}

	for request := range cacheRequests {
		if request.value != nil {
			cache[request.value.ID] = request.value
		} else if request.key != `` {
			request.value, _ = cache[request.key]
			ready(request)
		} else {
			request.list = make([]*performance, 0, len(cache))
			for _, perf := range cache {
				request.list = append(request.list, perf)
			}

			ready(request)
		}
	}
}
