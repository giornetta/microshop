package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   *ServerConfig   `yaml:"server"`
	Postgres *PostgresConfig `yaml:"postgres"`
	Kafka    *KafkaConfig    `yaml:"kafka"`
}

func FromYaml(filename string) (*Config, error) {
	yfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(yfile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSL      bool   `yaml:"ssl"`
}

func (c *PostgresConfig) ConnectionString() string {
	var ssl string
	if c.SSL {
		ssl = "enable"
	} else {
		ssl = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.Username, c.Password, c.Host, c.Port, c.Database, ssl)
}

type KafkaConfig struct {
	BrokerAddrs   []string `yaml:"brokers"`
	ConsumerGroup string   `yaml:"consumer-group"`
}
