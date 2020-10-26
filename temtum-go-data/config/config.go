package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	TX_OUTPUT_LIMIT      int    `env:"TX_OUTPUT_LIMIT" envDefault:"1000"`
	TX_MAX_VALIDITY_TIME int64  `env:"TX_MAX_VALIDITY_TIME" envDefault:"259200000"`
	ConStr               string `env:"ConStr"  envDefault:"user=su password=su dbname=testDb sslmode=disable"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &cfg
}
