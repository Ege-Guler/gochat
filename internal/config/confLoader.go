package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MQTT struct {
		Broker   string `yaml:"broker"`
		Port     int    `yaml:"port"`
		Topic    string `yaml:"topic"`
		ClientID string `yaml:"clientID"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		QoS      int    `yaml:"qos"`
		Retain   bool   `yaml:"retain"`
	} `yaml:"mqtt"`
	Chat struct {
		Username string `yaml:"username"`
	} `yaml:"chat"`
}

func (c *Config) LoadConf() error {

	f, err := os.ReadFile("conf.yaml")

	if err != nil {
		log.Fatal("missing config file, please check the config.yaml file", err)
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(f, c); err != nil {
		log.Fatal("Corrupted config file, please check the config.yaml file", err)
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return nil
}
