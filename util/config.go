package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configurations of the application.
// The values are read by viper from a config file or environment variable
type Config struct {
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	GrpcServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	HttpServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

// LoadConfig reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
