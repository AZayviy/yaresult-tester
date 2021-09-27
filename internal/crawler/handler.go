package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/AZayviy/yaresult-tester/internal/config"
	"github.com/AZayviy/yaresult-tester/internal/net"
	"github.com/AZayviy/yaresult-tester/internal/parser"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"
	"crypto/tls"
)

type httpHandler struct {
	crawler Crawler
}

type memoryStats struct {
	allocated      string
	totalAllocated string
	system         string
}

func getMemory() memoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return memoryStats{
		allocated:      fmt.Sprintf("%d MiB", m.Alloc/1024/1024),
		totalAllocated: fmt.Sprintf("%d MiB", m.TotalAlloc/1024/1024),
		system:         fmt.Sprintf("%d MiB", m.Sys/1024/1024),
	}
}

func NewHttpHandler(cfg config.CrawlerCfg) *httpHandler {
	redirectHandler := func(req *http.Request, via []*http.Request) error {
		return fmt.Errorf("redirect stopped")
	}

	transport := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,
        },
	}

	return &httpHandler{
		crawler: *NewCrawler(
			net.NewHttpClient(transport, cfg.RequestTimeout, redirectHandler),
			cfg,
		),
	}
}

func (hh *httpHandler) Handle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("Start memory usage: %+v\n", getMemory())
	query := r.URL.Query().Get("search")
	if query == "" {
		errMsg := "missing query argument `search`"
		log.Println(errMsg)
		writeErrorResponse([]byte(errMsg), w)
		return
	}

	res, err := hh.crawler.client.Get(fmt.Sprintf(parser.BaseYandexURL, url.QueryEscape(query)))
	if err != nil {
		log.Println(err)
		writeErrorResponse([]byte(err.Error()), w)
		return
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		writeErrorResponse([]byte(err.Error()), w)
		return
	}

	urls := parser.GetUrls(bytes)

	result := hh.crawler.ProcessUrlsInBatch(urls)

	writeJsonResponse(result, w)

	log.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	log.Printf("End memory usage: %+v\n", getMemory())
}

func writeJsonResponse(v interface{}, w http.ResponseWriter) {
	res, err := json.Marshal(v)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}

func writeErrorResponse(errMsg []byte, w http.ResponseWriter) {
	w.Write(errMsg)
}
