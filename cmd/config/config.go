/*
Copyright © 2026 German-Feskov
*/
package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Jira   JiraConfig   `mapstructure:"jira" validate:"required"`
	DB     DBConfig     `mapstructure:"db" validate:"required"`
	Broker BrokerConfig `mapstructure:"broker"`
	App    AppConfig    `mapstructure:"app" validate:"required"`
}

type JiraConfig struct {
	URL          string        `mapstructure:"url" validate:"required,url"`
	Token        string        `mapstructure:"token" validate:"required,min=10"`
	UnitWorkers  int           `mapstructure:"unit_workers" validate:"required,min=1,max=100"`
	MaxTimeRetry time.Duration `mapstructure:"max_time_retry" validate:"required,min=1ms"`
	MinTimeRetry time.Duration `mapstructure:"min_time_retry" validate:"required,min=1ms"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	User     string `mapstructure:"user" validate:"required,min=1"`
	Password string `mapstructure:"password" validate:"required,min=1"`
	DBName   string `mapstructure:"dbname" validate:"required,min=1"`
	SSLMode  string `mapstructure:"sslmode" validate:"required,oneof=disable require verify-ca verify-full"`
}

type BrokerConfig struct {
	URL           string `mapstructure:"url" validate:"omitempty,url"`
	Exchange      string `mapstructure:"exchange" validate:"omitempty,min=1"`
	ExchangeType  string `mapstructure:"exchange_type" validate:"omitempty,oneof=direct fanout topic headers"`
	Queue         string `mapstructure:"queue" validate:"omitempty,min=1"`
	RoutingKey    string `mapstructure:"routing_key" validate:"omitempty"`
	PrefetchCount int    `mapstructure:"prefetch_count" validate:"omitempty,min=1"`
	AutoReconnect bool   `mapstructure:"auto_reconnect"`
}

// TODO add some vars
type AppConfig struct {
	LogLevel string `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
}

//TODO add config for gRPC

var (
	appConfig *Config
	validate  *validator.Validate
)

func init() {
	validate = validator.New()
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := cfg.customValidate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	appConfig = &cfg
	return &cfg, nil
}

func GetConfig() *Config {
	if appConfig == nil {
		panic("config not loaded")
	}
	return appConfig
}

func (c *Config) customValidate() error {
	// Example of custom validation, TODO edit
	if c.Jira.Token == "{Your API token}" {
		return fmt.Errorf("jira token must be set")
	}

	if c.Jira.MinTimeRetry > c.Jira.MaxTimeRetry {
		return fmt.Errorf("jira min_time_retry must be less than max_time_retry")
	}

	return nil
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

func (j *JiraConfig) GetAPIURL(endpoint string) string {
	return fmt.Sprintf("%s/rest/api/2/%s", j.URL, endpoint)
}
