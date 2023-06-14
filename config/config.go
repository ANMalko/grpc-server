package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const configFile = "config.yml"

// type APIConfig struct {
// 	prefix string

// }

type Config struct {
	HTTPPort int `yaml:"http_port"`
	GRPCPort int `yaml:"grpc_port"`
}

func GetAPIConfig() *Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("File `%s` not found", configFile)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Unable parse `%s` file", configFile)
	}

	return &config
}
