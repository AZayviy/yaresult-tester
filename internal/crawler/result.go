package crawler

import "sync"

type singleItemAccess interface {
	Get(url string) (res int)
	Set(url string, res int)
}

type Urls struct {
	mu   sync.Mutex
	data map[string]int
}

func (u *Urls) Set(url string, res int) {
	u.mu.Lock()
	u.data[url] = res
	u.mu.Unlock()
}

func (u *Urls) Get(url string) (res int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	res = u.data[url]
	return
}

func NewUrlsMap(urls []string) *Urls {
	urlsMap := Urls{data: make(map[string]int)}
	for _, val := range urls {
		urlsMap.data[val] = 0
	}
	return &urlsMap
}
