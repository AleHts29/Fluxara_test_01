package config

import (
	"fluxara/internal/domain"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	GlobalConfig *Config
	once         sync.Once
	subsMu       sync.Mutex
	subscribers  []chan *Config
)

type Config struct {
	Server domain.Server `mapstructure:"server"`
	Db     domain.Db     `mapstructure:"db"`
}

func Load() {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath("./configs")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file\n")
		}

		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("Error config mapping\n")
		}

		GlobalConfig = &config
		fmt.Printf("Settings loaded successfully\n")
		broadcast(GlobalConfig)

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			var next Config
			if err := viper.Unmarshal(&next); err != nil {
				fmt.Printf("Error reloading configuration\n")
				return
			}

			GlobalConfig = &next
			broadcast(GlobalConfig)
		})
		fmt.Printf("Reload successful\n")
	})
}

func Get() *Config {
	return GlobalConfig
}

func (c *Config) Subscribe() <-chan *Config {
	ch := make(chan *Config, 1)
	subsMu.Lock()
	subscribers = append(subscribers, ch)
	subsMu.Unlock()

	if GlobalConfig != nil {
		ch <- GlobalConfig
	}
	return ch
}

func broadcast(cfg *Config) {
	subsMu.Lock()
	defer subsMu.Unlock()
	for _, ch := range subscribers {
		select {
		case ch <- cfg:
		default: // si está lleno, no bloqueo
		}
	}
}
