package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type JiraConfig struct {
	BaseURL string `yaml:"base_url"`
	Token   string `yaml:"token"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type AppConfig struct {
	Jira JiraConfig `yaml:"jira"`
	DB   DBConfig   `yaml:"db"`
	App  struct {
		LogLevel string `yaml:"log_level"`
	} `yaml:"app"`
}

func LoadDevConfig() (*AppConfig, error) {
	data, err := os.ReadFile("config/dev.yaml")
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
