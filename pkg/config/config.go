package config

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	EmptyPath = "EMPTY_PATH"
)

type ServiceConfig struct {
	Host string `mapstructure:"host"`
}

type CLIConfig struct {
	Mnemonic string `mapstructure:"mnemonic"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Environment   string        `mapstructure:"env"`
	ServiceConfig ServiceConfig `mapstructure:"service_config"`
	CLIConfig     CLIConfig     `mapstructure:"cli_config"`
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}

	// Get config path from flags, if flags not exist get from env
	path := pflag.String("app-cfg", EmptyPath, "config path")
	mnemonic := pflag.String("mnemonic", "", "mnemonic phrase for wallet authorization")
	password := pflag.String("password", "", "password for wallet")
	pflag.Parse()

	if *path == EmptyPath {
		viper.AutomaticEnv()
		*path = viper.Get("APP_CFG_PATH").(string)
	}

	v := viper.New()
	v.SetConfigFile(*path)
	slog.Info("Read config", "path:", *path)
	if err := v.ReadInConfig(); err != nil {
		return nil, errors.New("failed to get yaml config")
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.New("failed to pull config")
	}

	// Override CLI config with command line flags
	if *mnemonic != "" {
		cfg.CLIConfig.Mnemonic = *mnemonic
	}
	if *password != "" {
		cfg.CLIConfig.Password = *password
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	slog.Info("config", "cfg:", cfg)
	return cfg, nil
}
