package config

import (
	"github.com/spf13/viper"
)

const (
	// EnvDev const represents dev environment
	EnvDev = "dev"
	// EnvStaging const represents staging environment
	EnvStaging = "staging"
	// EnvProduction const represents production environment
	EnvProduction = "production"
)

// Initialize ...
func Initialize() {

	viper.AutomaticEnv()

	viper.SetDefault("ENV", EnvDev)
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASS", "SecretPassword")
	viper.SetDefault("DB_NAME", "slotracker")

}
