package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Jira   JiraConfig   `mapstructure:"jira"   validate:"required"`
	DB     DBConfig     `mapstructure:"db"     validate:"required"`
	Broker BrokerConfig `mapstructure:"broker" validate:"required"`
	App    AppConfig    `mapstructure:"app"    validate:"required"`
}

type JiraConfig struct {
	URL          string        `mapstructure:"url"            validate:"required,url"`
	Token        string        `mapstructure:"token"          validate:"required,min=1"`
	UnitWorkers  int           `mapstructure:"unit_workers"   validate:"required,min=1"`
	MaxTimeRetry time.Duration `mapstructure:"max_time_retry" validate:"required,gt=0"`
	MinTimeRetry time.Duration `mapstructure:"min_time_retry" validate:"required,gt=0"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"     validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port"     validate:"required,min=1,max=65535"`
	User     string `mapstructure:"user"     validate:"required,min=1"`
	Password string `mapstructure:"password" validate:"required,min=1"`
	DBName   string `mapstructure:"dbname"   validate:"required,min=1"`
	SSLMode  string `mapstructure:"sslmode"  validate:"required,oneof=disable require verify-ca verify-full"`
}

type BrokerConfig struct {
	URL           string `mapstructure:"url"            validate:"required,url"`
	Exchange      string `mapstructure:"exchange"       validate:"required,min=1"`
	ExchangeType  string `mapstructure:"exchange_type"  validate:"required,oneof=direct topic fanout headers"`
	Queue         string `mapstructure:"queue"          validate:"required,min=1"`
	RoutingKey    string `mapstructure:"routing_key"    validate:"required,min=1"`
	PrefetchCount int    `mapstructure:"prefetch_count" validate:"required,min=1"`
	AutoReconnect bool   `mapstructure:"auto_reconnect"`
}

type AppConfig struct {
	LogLevel string `mapstructure:"log_level" validate:"required,oneof=debug info warn error"`
}

// nolint: gochecknoglobals
var (
	appConfig *Config
	validate  *validator.Validate
)

// nolint: gochecknoinits
func init() {
	validate = validator.New()
}

func LoadConfig() (*Config, error) {
	var cfg Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	err = validate.Struct(&cfg)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = cfg.customValidate()
	if err != nil {
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
	return nil
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}
