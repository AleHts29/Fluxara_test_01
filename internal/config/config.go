package config

import (
	"fluxara/internal/domain"
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	GlobalConfig *Config
	once         sync.Once
	subsMu       sync.Mutex
	subscribers  []chan *Config
)

type Config struct {
	Server   domain.Server `mapstructure:"server"`
	Db       domain.Db     `mapstructure:"db"`
	DbGergal domain.Db     `mapstructure:"dbGergal"`
}

func Load() {
	once.Do(func() {
		// 🔹 1) Cargar .env al entorno del sistema
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("No se pudo cargar .env (se usan env del sistema)")
		}

		viper.AutomaticEnv()

		cfg := Config{
			Server: domain.Server{
				Host: viper.GetString("SERVER_HOST"),
				Port: viper.GetString("SERVER_PORT"),
				TLS: domain.TLSServerConf{
					Enable:              viper.GetBool("SERVER_TLS_ENABLE"),
					HttpsPort:           viper.GetString("SERVER_TLS_HTTPS_PORT"),
					CertFile:            viper.GetString("SERVER_TLS_CERT_FILE"),
					KeyFile:             viper.GetString("SERVER_TLS_KEY_FILE"),
					RedirectHTTPToHTTPS: viper.GetBool("SERVER_TLS_REDIRECT_HTTP_TO_HTTPS"),
				},
			},
			Db: domain.Db{
				Connection:  viper.GetString("DB_CONNECTION"),
				Host:        viper.GetString("DB_HOST"),
				Port:        viper.GetString("DB_PORT"),
				User:        viper.GetString("DB_USER"),
				Password:    viper.GetString("DB_PASSWORD"),
				Name:        viper.GetString("DB_NAME"),
				SslMode:     viper.GetString("DB_SSLMODE"),
				Retries:     viper.GetInt("DB_RETRIES"),
				TimeRetries: viper.GetInt("DB_TIME_RETRIES"),
			},
		}

		fmt.Printf("ESTO ES CFG DESDE CONFIG: %+v\n\n", &cfg)

		GlobalConfig = &cfg
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
