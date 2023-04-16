package config

import (
	"strings"

	"github.com/spf13/viper"
)

var Configuration GlobalConfiguration

// Global Config
type GlobalConfiguration struct {
	Database DatabaseConfig `mapstructure:"database"`
	API      OpenAPIConfig  `mapstructure:"openapi"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int32  `mapstructure:"port"`
	Db       string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type OpenAPIConfig struct {
	URL string `mapstructure:"url"`
}

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("configuration")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// Default value for DB port
	viper.SetDefault("database.port", 5432)

	// Read config
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Configuration)
	if err != nil {
		panic(err)
	}
}
