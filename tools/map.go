package tools

// ConcurrentMap is a map[string]MapContent with concurrent access
type ConcurrentMap struct {
	req     chan request
	content map[string]MapContent
}

// MapContent stored in map
type MapContent interface {
	GetID() string
}

type request struct {
	action    string
	key       string
	ioContent chan MapContent
}

// Get entry in map
func (c *ConcurrentMap) Get(key string) MapContent {
	request := request{action: `get`, key: key, ioContent: make(chan MapContent)}
	c.req <- request

	return <-request.ioContent
}

// Push entry in map
func (c *ConcurrentMap) Push(entry MapContent) {
	request := request{action: `push`, ioContent: make(chan MapContent)}
	c.req <- request

	request.ioContent <- entry
	close(request.ioContent)
}

// List map content to output channel
func (c *ConcurrentMap) List() <-chan MapContent {
	results := make(chan MapContent, len(c.content))
	c.req <- request{action: `list`, ioContent: results}

	return results
}

// CreateConcurrentMap in a subroutine
func CreateConcurrentMap(contentSize int, channelSize int) *ConcurrentMap {
	concurrentMap := ConcurrentMap{req: make(chan request, channelSize), content: make(map[string]MapContent, contentSize)}

	go func() {
		for request := range concurrentMap.req {
			if request.action == `list` {
				for _, perf := range concurrentMap.content {
					request.ioContent <- perf
				}
				close(request.ioContent)
			} else if request.action == `push` {
				entry := <-request.ioContent
				concurrentMap.content[entry.GetID()] = entry
			} else if request.action == `get` {
				if entry, ok := concurrentMap.content[request.key]; ok {
					request.ioContent <- entry
				}
				close(request.ioContent)
			}
		}
	}()

	return &concurrentMap
}
