package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig
	DB     DBConfig
	JWT    JWTConfig
	SMTP   SMTPConfig
	Upload UploadConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port int
}

type DBConfig struct {
	URI  string
	Name string
}

type JWTConfig struct {
	AccessSecret        string `mapstructure:"access_secret"`
	RefreshSecret       string `mapstructure:"refresh_secret"`
	AccessExpiryMinutes int    `mapstructure:"access_expiry_minutes"`
	RefreshExpiryDays   int    `mapstructure:"refresh_expiry_days"`
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type UploadConfig struct {
	Dir       string
	MaxSizeMB int `mapstructure:"max_size_mb"`
}

func Load() (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// ENV vars override yaml values
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicit env bindings
	viper.BindEnv("db.uri", "MONGO_URI")
	viper.BindEnv("db.name", "MONGO_DB_NAME")
	viper.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	viper.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	viper.BindEnv("app.port", "APP_PORT")
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.username", "SMTP_USERNAME")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")
	viper.BindEnv("smtp.from", "SMTP_FROM")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
