package config

import (
	"fmt"
	"time"
	"github.com/caarlos0/env/v6"
)

type HttpCfg struct {
	Host string `env:"HTTP_HOST"`
	Port uint16 `env:"SERVER_PORT" envDefault:"8999"`
}

type CrawlerCfg struct {
	ResponseTimeout time.Duration `env:"RESPONSE_TIMEOUT" envDefault:"30s"`
	StartThreads    int           `env:"START_THREADS" envDefault:"5"`
	Iterations      int           `env:"ITERATIONS_LIMIT" envDefault:"100"`
	RequestTimeout  time.Duration `env:"REQUEST_TIMEOUT" envDefault:"3s"`
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
