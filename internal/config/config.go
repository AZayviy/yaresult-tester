package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type HttpCfg struct {
	Host string `env:"HTTP_HOST"`
	Port uint16 `env:"SERVER_PORT" envDefault:"8999"`
}

type CrawlerCfg struct {
	ResponseTimeout time.Duration `env:"SERVER_RESPONSE_TIMEOUT" envDefault:"29s"`
	StartThreads    int           `env:"CRAWLER_START_THREADS" envDefault:"5"`
	Iterations      int           `env:"CRAWLER_ITERATIONS_LIMIT_PER_URL" envDefault:"10"`
	RequestTimeout  time.Duration `env:"CRAWLER_REQUEST_TIMEOUT" envDefault:"3s"`
}

type Cfg struct {
	AppEnv  string `env:"APP_ENV" envDefault:"dev"`
	Http    HttpCfg
	Crawler CrawlerCfg
}

func Init() *Cfg {
	cfg := new(Cfg)
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Sprintf("can't get config: %s", err))
	}
	return cfg
}
