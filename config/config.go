package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres" envPrefix:"POSTGRES_"`
	Kafka    KafkaConfig    `yaml:"kafka" envPrefix:"KAFKA_"`
}

func FromYaml(filename string) (*Config, error) {
	yfile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(yfile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func FromEnv() (*Config, error) {
	var cfg Config

	opts := env.Options{UseFieldNameByDefault: true}

	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Port     int    `yaml:"port"`
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
	ConsumerGroup string   `yaml:"consumer-group" env:"CG"`
	BrokerAddrs   []string `yaml:"brokers" env:"BROKERS"`
}
