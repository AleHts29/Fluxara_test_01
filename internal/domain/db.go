package domain

import (
	"time"
)

type Db struct {
	Connection  string `mapstructure:"connection"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	SslMode     string `mapstructure:"sslmode"`
	Retries     int    `mapstructure:"retries"`
	TimeRetries int    `mapstructure:"time_retries"`
}

type Product struct {
	ID         string
	SKU        string
	Name       string
	Category   string
	PriceCents int32
	Stock      int32
	IsActive   bool
	CreatedAt  time.Time
}
