package config

import (
	"encoding/json"
	"net"
	"os"
)

type LoggerConfig struct {
	Level   string `json:"level"`
	LogFile string `json:"logfile"`
}

type EndpointConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (c EndpointConfig) GetEndpoint() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type StorageConfig struct {
	UseDb    bool   `json:"useDb"`
	Host     string `json:"host"`
	Port     int16  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	LoggerConfig  LoggerConfig   `json:"logger"`
	HttpConfig    EndpointConfig `json:"http"`
	GrpcConfig    EndpointConfig `json:"grpc"`
	StorageConfig StorageConfig  `json:"storage"`
}

func NewConfig(filename string) (Config, error) {
	var cfg Config

	arr, err := os.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(arr, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
