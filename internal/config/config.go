package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Logger   LoggerConfig   `mapstructure:",squash"`
}

type AppConfig struct {
	Name                 string `mapstructure:"app_name"`
	Env                  string `mapstructure:"app_env"`
	HTTPPort             int    `mapstructure:"app_http_port"`
	SecretKey            string `mapstructure:"app_secret_key"`
	RestaurantServiceUrl string `mapstructure:"restaurant_service_url"`
}
 
type DatabaseConfig struct {
	Host     string `mapstructure:"database_host"`
	Port     int    `mapstructure:"database_port"`
	User     string `mapstructure:"database_user"`
	Password string `mapstructure:"database_password"`
	DBName   string `mapstructure:"database_dbname"`
}

type LoggerConfig struct {
	Level string `mapstructure:"logger_level"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Since keys are flat, we don't need a replacer for .env, but for env variables like APP_NAME,
	// viper automatically looks for uppercase versions if AutomaticEnv is enabled.
	// Actually, viper will lowercase all env vars, so APP_NAME becomes app_name.

	_ = viper.ReadInConfig() // Ignore error if .env file is missing, fallback to env vars

	// Set defaults
	viper.SetDefault("app_name", "gopher-foody-order-service")
	viper.SetDefault("app_env", "development")
	viper.SetDefault("app_http_port", 8082)
	viper.SetDefault("restaurant_service_url", "http://localhost:8002")
	viper.SetDefault("logger_level", "debug")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
