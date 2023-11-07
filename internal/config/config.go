package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	URLAlphabet string     `yaml:"url_alphabet" env-required:"true"`
	URLLength   int        `yaml:"url_length" env-default:"10"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	DB          DB         `yaml:"db"`
}

type HTTPServer struct {
	StorageType string
	Host        string `yaml:"host" env-default:"localhost"`
	Port        string `yaml:"port" env-default:"8000"`
}

type DB struct {
	Host       string `yaml:"host" env-required:"true"`
	Port       string `yaml:"port" env-required:"true"`
	User       string
	DBName     string
	DBPassword string
}

func Load() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	dbName := os.Getenv("POSTGRES_DB")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbUser := os.Getenv("POSTGRES_USER")
	storageType := os.Getenv("STORAGE_TYPE")

	switch {
	case storageType == "":
		log.Fatal("STORAGE_PATH not found")
	case configPath == "":
		log.Fatal("CONFIG_PATH not found")
	case dbName == "":
		log.Fatal("POSTGRES_DB not found")
	case dbUser == "":
		log.Fatal("POSTGRES_USER not found")
	case dbPassword == "":
		log.Fatal("POSTGRES_PASSWORD not found")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	cfg.DB.DBName = dbName
	cfg.DB.DBPassword = dbPassword
	cfg.DB.User = dbUser
	cfg.HTTPServer.StorageType = storageType

	return &cfg
}
