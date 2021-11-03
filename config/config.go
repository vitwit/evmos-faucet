package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	UI struct {
		// port for running the ui network
		Port uint64 `yaml:"port" mapstructure:"port" validate:"required"`
	}

	Faucet struct {
		// account_prefix for bech32 encoding
		AccountPrefix string `yaml:"account_prefix" mapstructure:"account_prefix" validate:"required"`
		// env_prefix for storing the all env variables
		EnvPrefix string `yaml:"env_prefix" mapstructure:"env_prefix"`
		// amount is tokens per request
		Amount int64 `yaml:"amount" mapstructure:"amount" validate:"required"`
		// maximum tokens allowed for an account
		MaxTokens int64 `yaml:"max_tokens" mapstructure:"max_tokens" validate:"required"`
		// tendermint node address
		Node string `yaml:"node" mapstructure:"node" validate:"required"`
		// Lcd for quering the account balances
		Lcd string `yaml:"lcd" mapstructure:"lcd" validate:"required"`
		// chain denom
		Denom string `yaml:"denom" mapstructure:"denom" validate:"required"`
		// Decimals
		Decimals int `yaml:"decimals" json:"decimals" mapstructure:"decimals"  validate:"required"`
	}
}

//ReadConfigFromFile to read config details using viper
func ReadConfigFromFile() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config/")
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error while reading config.toml: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config.toml to application config: %v", err)
	}

	// validating the config
	if err := cfg.Validate(); err != nil {
		log.Fatalf("error occurred in config validation: %v", err)
	}

	return &cfg, nil
}

//Validate config struct
func (c *Config) Validate(e ...string) error {
	v := validator.New()
	if len(e) == 0 {
		return v.Struct(c)
	}
	return v.StructExcept(c, e...)
}
