package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Database   `yaml:"database"`
	HTTPServer `yaml:"http_server"`
}

type Database struct {
    DbURL  string `yaml:"dburl" env-default:"mongodb://localhost:27017"`
    DbName string `yaml:"dbname" env-default:"taskdb"`
}

type HTTPServer struct {
    Address      string `yaml:"address" env-default:"localhost:50051"`
    KafkaTopic   string `yaml:"kafka_topic" env-default:"default_topic"`
    KafkaGroupID string `yaml:"kafka_group_id" env-default:"default_group"`
}

func MustLoad() *Config {
	configPath := "../config/config.yml"
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatalf("cannot parse YAML config: %s", err)
	}

	return &cfg
}
