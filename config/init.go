package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ProfilerStatus     string `default:"true"`
	GroupId            string `default:"ANLYTIC_GROUP_ID"`
	DbConnectionString string `envconfig:"DB_CONNECTION_STRING"`
	DbType             string `envconfig:"DB" default:"PG"`
	LogLevel           string `envconfig:"LOGLEVEL" default:"debug"`
	AuthAddress        string `envconfig:"AUTHSERVICE_ADDRESS" default:"localhost:8015"`
	GrpcAddress        string `envconfig:"GPRC_ADDRESS" default:"localhost:4000"`

	StateKafkaTopic  string   `default:"team21-TASK.STATE.CHANGED"`
	EventKafkaTopic  string   `default:"team21-TASK.APPROVE.SET"`
	NotifyKafkaTopic string   `default:"team21-NOTIFICATION.SEND.SUCCESS"`
	Brokers          []string `envconfig:"KAFKA_BROKERS" default:"91.185.95.87:9094"`
}

func Init() (*Config, error) {
	var config Config
	err := envconfig.Process("ANALYTIC", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil

}
