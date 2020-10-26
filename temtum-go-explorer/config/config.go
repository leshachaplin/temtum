package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerPort         int    `env:"Port" envDefault:"1234"`
	TopicBlocks        string `env:"Topic" envDefault:"blocks168"`
	KafkaUrl           string `env:"Port" envDefault:"localhost:9092"`
	BatchReadGroup     string `env:"BatchReadGroup" envDefault:"batch"`
	Origin             string `env:"Origin" envDefault:"http://localhost:6774/"`
	Url                string `env:"Url" envDefault:"ws://localhost:6774/message"`
	Host               string `env:"Url" envDefault:"localhost:6774"`
	MsgPath            string `env:"Url" envDefault:"/message"`
	DbDriver           string `env:"DbDriver" envDefault:"postgres"`
	ConStr             string `env:"ConStr"  envDefault:"user=su password=su dbname=explorer sslmode=disable"`
	BatchSize          int    `env:"BatchSize" envDefault:"2000000000"`
	MaxBytesPerMessage int    `env:"MaxBytesPerMessage" envDefault:"1000000"`
	MaxBlocksAmount    int    `env:"MaxBlocksAmount" envDefault:"10000"`
	Timeout            int    `env:"Timeout" envDefault:"10"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &cfg
}
