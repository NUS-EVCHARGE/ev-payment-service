package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	GetUserUrl       = ""
	GetUserUrlPrefix = "/user/get_user_info"
	GetBookingUrl    = ""
	GetBookingPrefix = "/booking"
)

type Config struct {
	Dsn            string `mapstructure:"dsn"`
	HttpAddress    string `mapstructure:"http_address"`
	BaseServiceUrl string `mapstructure:"base_service_url"`
	MongoDBURL     string `mapstructure:"MongoDB_url"`
	StripeKey      string `mapstructure:"stripe_key"`
}

func ParseConfig(filename string) (*Config, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename cannout be empty")
	}
	var config *Config
	v := viper.New()
	v.SetConfigFile(filename)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	GetUserUrl = config.BaseServiceUrl + GetUserUrlPrefix
	GetBookingUrl = config.BaseServiceUrl + GetBookingPrefix
	return config, nil
}
