package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server Server `yaml:"server"`
	Mongodb MongoDB `yaml:"mongodb"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type MongoDB struct {
	Uri string `yaml:"uri"`
}

func LoadConf() (*Config, error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	viper.SetDefault("server", map[string]string{"host": "127.0.0.1", "port": "9527"})
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}