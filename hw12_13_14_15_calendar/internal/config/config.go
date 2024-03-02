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

type RabbitConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Queue    string `json:"queue"`
}

type TasksConfig struct {
	CleanupIntervalMin     int32 `json:"cleanup_interval_min"`
	CleanupPeriodDays      int32 `json:"cleanup_period_days"`
	NotifyCheckIntervalMin int32 `json:"notify_check_interval_min"`
}

type CalendarConfig struct {
	LoggerConfig  LoggerConfig   `json:"logger"`
	HttpConfig    EndpointConfig `json:"http"`
	GrpcConfig    EndpointConfig `json:"grpc"`
	StorageConfig StorageConfig  `json:"storage"`
}

type SchedulerConfig struct {
	LoggerConfig  LoggerConfig  `json:"logger"`
	StorageConfig StorageConfig `json:"storage"`
	RabbitConfig  RabbitConfig  `json:"rabbit"`
	TasksConfig   TasksConfig   `json:"tasks"`
}

type SenderConfig struct {
	LoggerConfig LoggerConfig `json:"logger"`
	RabbitConfig RabbitConfig `json:"rabbit"`
}

func NewCalendarConfig(filename string) (CalendarConfig, error) {
	var cfg CalendarConfig
	return cfg, newConfig(filename, &cfg)
}

func NewSchedulerConfig(filename string) (SchedulerConfig, error) {
	var cfg SchedulerConfig
	return cfg, newConfig(filename, &cfg)
}

func NewSenderConfig(filename string) (SenderConfig, error) {
	var cfg SenderConfig
	return cfg, newConfig(filename, &cfg)
}

func newConfig(filename string, dest interface{}) error {
	arr, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(arr, dest)
	if err != nil {
		return err
	}

	return nil
}
