package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"` // hours
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.JWT.Expire == 0 {
		cfg.JWT.Expire = 72
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "upp-secret-key"
	}
	return &cfg, nil
}

func (d *DatabaseConfig) DSN() string {
	return "host=" + d.Host + " port=" + d.Port + " user=" + d.User + " password=" + d.Password + " dbname=" + d.DBName + " sslmode=" + d.SSLMode
}
