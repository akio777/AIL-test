package db

import (
	"sync"

	"github.com/uptrace/bun"
)

type Config struct {
	Host      string `required:"true"`
	Port      int    `required:"true"`
	User      string `required:"true"`
	Password  string `required:"true"`
	Name      string `required:"true"`
	PoolSize  int    `split_words:"true"`
	AppName   string `split_words:"true"`
	SSLEnable string `split_words:"true" envconfig:"SSL_ENABLE"`

	once       sync.Once
	connection *bun.DB
}

type Driver string

const (
	POSTGRES = Driver("postgres")
)
