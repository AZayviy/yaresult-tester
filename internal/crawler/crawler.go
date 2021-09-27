package crawler

import (
	"github.com/AZayviy/yaresult-tester/internal/config"
	"log"
	"net/http"
	"time"
)

type Crawler struct {
	client          *http.Client
	iterationsLimit int
	responseTimeout time.Duration
	requestTimeout  time.Duration
	startThreads    int
}

func NewCrawler(client *http.Client, cfg config.CrawlerCfg) *Crawler {
	return &Crawler{
		client:          client,
		iterationsLimit: cfg.Iterations,
		responseTimeout: cfg.ResponseTimeout,
		startThreads:    cfg.StartThreads,
		requestTimeout:  cfg.RequestTimeout,
	}
}

//ProcessUrlsInBatch test each url and returns map url=>recommended concurrent requests
func (c *Crawler) ProcessUrlsInBatch(urls []string) (res map[string]int) {
	exitCh := make(chan bool)
	syncUrls := NewUrlsMap(urls)

	for _, url := range urls {
		go c.proccessSingleUrl(url, syncUrls, exitCh)
	}

	timeout := time.After(c.responseTimeout - c.requestTimeout)

	for i := 0; i < len(urls); i++ {
		select {
		case <-exitCh:
			continue
		case <-timeout:
			log.Printf("timeout\n")
			res = make(map[string]int)
			for _, url := range urls {
				res[url] = syncUrls.Get(url)
			}
			return
		}
	}

	return syncUrls.data
}

func (c *Crawler) proccessSingleUrl(url string, syncUrls *Urls, exitCh chan<- bool) {
	current := c.startThreads
	previousUpscaled := false
	for i := 0; i < c.iterationsLimit; i++ {
		res := c.testUrl(url, current)
		step := current / 2
		if res {
			syncUrls.Set(url, current)
			current = current + step
			previousUpscaled = true
		} else {
			if previousUpscaled {
				break
			}
			current = current - step
		}
		if current <= 1 {
			syncUrls.Set(url, 1)
			break
		}
	}
	exitCh <- true
}

func (c *Crawler) testUrl(url string, limit int) bool {
	ch := make(chan bool)
	for i := 0; i < limit; i++ {
		go c.makeRequest(url, ch)
	}
	res := true

	for i := 0; i < limit; i++ {
		accessible := <-ch
		res = res && accessible
	}
	return res
}

func (c *Crawler) makeRequest(url string, ch chan<- bool) {
	resp, err := c.client.Get(url)
	if err != nil {
		ch <- false
		return
	}
	defer resp.Body.Close()
	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		ch <- false
		return
	}
	ch <- true
}
