package main

import (
	"fmt"
	"github.com/AZayviy/yaresult-tester/internal/config"
	"github.com/AZayviy/yaresult-tester/internal/crawler"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	cfg := config.Init()
	handler := crawler.NewHttpHandler(cfg.Crawler)
	http.HandleFunc("/sites", handler.Handle)

	addr := fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
