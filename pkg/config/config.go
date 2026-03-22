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

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DBHosts struct {
	DBHost string `mapstructure:"host"`
	DBPort int    `mapstructure:"port"`
}

type DBConfig struct {
	DBUser         string  `mapstructure:"user"`
	DBName         string  `mapstructure:"name"`
	DBPassword     string  `mapstructure:"password"`
	DBHosts        DBHosts `mapstructure:"hosts"`
	MigrationsPath string  `mapstructure:"migrations_path"`
}

type JWT struct {
	SecretKey string `mapstructure:"secret_key"`
}

type ServiceConfig struct {
	Host string `mapstructure:"host"`
}

type Config struct {
	Environment   string        `mapstructure:"env"`
	Line          int64         `mapstructure:"line"`
	ServerConfig  ServerConfig  `mapstructure:"server_config"`
	DBConfig      DBConfig      `mapstructure:"db_config"`
	ServiceConfig ServiceConfig `mapstructure:"service_config"`
	JWT           JWT           `mapstructure:"jwt"`
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}

	// Get config path from flags, if flags not exist get from env
	path := pflag.String("app-cfg", EmptyPath, "config path")
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

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// JWT
	cfg.JWT.SecretKey = v.GetString("jwt.secret_key")

	// DB PASSWORD
	cfg.DBConfig.DBPassword = v.GetString("db_password")

	slog.Info("config", "cfg:", cfg)
	return cfg, nil
}
