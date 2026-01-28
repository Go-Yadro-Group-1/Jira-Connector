package config

type Config struct {
	Database struct {
		DSN string
	}
}

func Load() (*Config, error) {
	return &Config{}, nil
}
