package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ZapPresetExample     = "example"
	ZapPresetDevelopment = "development"
)

type DBConfig struct {
	DSN string
}

type ZapConfig struct {
	Preset string
}

type TracerConfig struct {
	Host            string
	Port            string
	Service         string
	Env             string
	Version         string
	EnableMockAgent bool
	EnableDebugMode bool
}

type Config struct {
	DB     DBConfig
	Zap    ZapConfig
	Tracer TracerConfig
}

func LoadConfig() (Config, error) {
	var c Config

	if err := viper.ReadInConfig(); err != nil {
		return c, fmt.Errorf("could not read config: %w", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return c, nil
}

func InitProfileConfig(path string, fileType string) {
	profile := os.Getenv("PROFILE")
	if profile != "" {
		viper.SetConfigName(profile)
	} else {
		viper.SetConfigName("app")
	}

	viper.AddConfigPath(filepath.Join(path, "config"))
	viper.SetConfigType(fileType)
	viper.AutomaticEnv()
}
