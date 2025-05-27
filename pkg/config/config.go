package config

import (
	"github.com/caarlos0/env/v11"

	bugsnagAPI "github.com/sazap10/bugsnag-api-go"
)

// Config holds the configuration for the application.
type Config struct {
	// bugsnag auth token
	AuthToken string `env:"BUGSNAG_AUTH_TOKEN,required"`
	// bugsnag endpoint
	Endpoint string `env:"BUGSNAG_ENDPOINT" envDefault:"https://api.bugsnag.com"`

	APIClient *bugsnagAPI.Client
}

// NewConfig creates a new Config struct and populates it with environment variables.
// It returns an error if any required environment variables are missing or if there is an error parsing them.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	cfg.APIClient = bugsnagAPI.NewClient(
		cfg.AuthToken,
		bugsnagAPI.WithBaseURL(cfg.Endpoint),
	)
	return cfg, nil
}
