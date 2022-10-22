package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"false"`
	Listen        struct {
		SocketFile string `env:"SOCKET_FILE" env-default:"app.sock"`
		Type       string `env:"LISTEN_TYPE" env-default:"port"`
		BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"PORT" env-default:"10000"`
	}
	Storage   StorageConfig `env:"STORAGE"`
	AppConfig struct {
		LogLevel string `env:"LOG_LEVEL" env-default:"trace"`
	}
	Kafka struct {
		URL   string `env:"KAFKA_URL" env-default:"kafka:9092"`
		Topic string `env:"TOPIC" env-default:"proxy_checker"`
	}
}

type StorageConfig struct {
	Host     string `env:"HOST" env-default:"postgres"`
	Port     string `env:"PORT" env-default:"5432"`
	Database string `env:"DATABASE" env-default:"postgres"`
	Username string `env:"POSTGRES_USERNAME" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Printf("gather config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "An error occurred during reading config"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
