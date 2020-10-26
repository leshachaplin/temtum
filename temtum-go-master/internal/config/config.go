package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerPort         int    `env:"Port" envDefault:"6774"`
	ClientPort         int    `env:"Port" envDefault:"8668"`
	TopicTransactions  string `env:"Topic" envDefault:"transactions94"`
	TopicBlocks        string `env:"Topic" envDefault:"blocks199"`
	KafkaUrl           string `env:"Port" envDefault:"localhost:9092"`
	TransactionGroup   string `env:"TransactionGroup" envDefault:"transaction"`
	BlockGroup         string `env:"BlockGroup" envDefault:"block"`
	Origin             string `env:"Origin" envDefault:"http://localhost:8668/"`
	Url                string `env:"Url" envDefault:"ws://localhost:8668/role"`
	ExitUrl            string `env:"ExitUrl" envDefault:"ws://localhost:8668/exit"`
	DbDriver           string `env:"DbDriver" envDefault:"postgres"`
	ConStr             string `env:"ConStr"  envDefault:"user=su password=su dbname=master sslmode=disable"`
	MaxBytesPerMessage int    `env:"MaxBytesPerMessage" envDefault:"2000000000"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &cfg
}
