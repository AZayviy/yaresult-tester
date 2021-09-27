package crawler

import (
	"github.com/AZayviy/yaresult-tester/internal/config"
	"net/http"
	"time"
	"fmt"
)

type Crawler struct {
	urls   *Urls
	client http.Client
	cfg    config.CrawlerCfg
}

func NewCrawler(urls *Urls, client http.Client, cfg config.CrawlerCfg) *Crawler {
	return &Crawler{urls: urls, client: client, cfg: cfg}
}

func (c *Crawler) ProcessUrls(urls []string) (res map[string]int) {
	exitCh := make(chan bool)

	for _, url := range urls {
		go c.proccessUrl(url, exitCh)
	}

	for i := 0; i < len(urls); i++ {
		select {
		case <-exitCh:
			continue
		case <-time.After(c.cfg.ResponseTimeout * time.Second):
			res = make(map[string]int)
			for _, url := range urls {
				res[url] = c.urls.Get(url)
			}
			return
		}
	}

	return c.urls.data
}

func (c *Crawler) proccessUrl(url string, exitCh chan<- bool) {
	current := c.cfg.StartThreads
	previousUpscaled := false
	for i := 0; i < c.cfg.Iterations; i++ {
		res := c.testUrl(url, current)
		step := current / 2
		if res {
			c.urls.Set(url, current)
			current = current + step
			previousUpscaled = true
		} else {
			if previousUpscaled {
				break
			}
			current = current - step
		}
		if current <= 1 {
			c.urls.Set(url, 1)
			break
		}
	}
	fmt.Printf("Tested %s result: %d\n", url, c.urls.Get(url))
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
