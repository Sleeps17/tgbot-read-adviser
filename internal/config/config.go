package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env"`
	LogsPath       string `yaml:"logs_path"`
	BotConfig      `yaml:"bot_config"`
	DataBaseConfig `yaml:"data_base_config"`
}

type BotConfig struct {
	TgBotHost  string `yaml:"tg_bot_host"`
	TgBotToken string `yaml:"tg_bot_token"`
	BatchSize  int    `yaml:"batch_size"`
}

type DataBaseConfig interface{}

type SQLiteConfig struct {
	StoragePath string `yaml:"storage_path"`
}

type PostgresqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// database_name: tgbot
// user_name: tgbot_user
// password: Pavel19122004

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error with reading config file: %v", err)
	}

	return &cfg
}
