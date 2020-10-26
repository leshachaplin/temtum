package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ClientPort              int    `env:"Port" envDefault:"8668"`
	Network                 string `env:"Network" envDefault:"tcp"`
	TopicTransactions       string `env:"Topic" envDefault:"transactions94"`
	TopicBlocks             string `env:"Topic" envDefault:"blocks199"`
	KafkaUrl                string `env:"Port" envDefault:"localhost:9092"`
	TransactionGroup        string `env:"TransactionGroup" envDefault:"transaction"`
	BatchBlockGroup         string `env:"BatchReadGroup" envDefault:"batch"`
	BlockGroup              string `env:"BlockGroup" envDefault:"block"`
	Origin                  string `env:"Origin" envDefault:"http://localhost:6774/"`
	Host                    string `env:"Host" envDefault:"localhost:6774"`
	DbDriver                string `env:"DbDriver" envDefault:"postgres"`
	RolePath                string `env:"RolePath" envDefault:"/role"`
	MsgPath                 string `env:"MsgPath" envDefault:"/message"`
	ExitPath                string `env:"ExitPath" envDefault:"/exit"`
	SyncPath                string `env:"SyncPath" envDefault:"/sync"`
	ConStrExp               string `env:"ConStrExp"  envDefault:"user=su password=su dbname=explorer sslmode=disable"`
	ConStr                  string `env:"ConStr"  envDefault:"user=su password=su dbname=blocks sslmode=disable"`
	MaxBytesPerMessage      int    `env:"MaxBytesPerMessage" envDefault:"1000000"`
	BatchSize               int    `env:"BatchSize" envDefault:"1000000000"`
	MaxBlocksAmount         int    `env:"MaxBlocksAmount" envDefault:"10000"`
	MaxTransactionsAmount   int    `env:"MaxTransactionsAmount" envDefault:"10000"`
	BatchTimeout            int    `env:"BatchTimeout" envDefault:"10"`
	ReadTransactionsTimeout int    `env:"ReadTransactionsTimeout" envDefault:"12"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &cfg
}
