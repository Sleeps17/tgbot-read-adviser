package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	LogsPath    string `yaml:"logs_path"`
	BotConfig   `yaml:"bot_config"`
}

type BotConfig struct {
	TgBotHost  string `yaml:"tg_bot_host"`
	TgBotToken string `yaml:"tg_bot_token"`
	BatchSize  int    `yaml:"batch_size"`
}

func MustLoad() *Config {
	config_path := os.Getenv("CONFIG_PATH")
	if config_path == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(config_path, &cfg); err != nil {
		log.Fatalf("Error with reading config file: %v", err)
	}

	return &cfg
}
