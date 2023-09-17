package server

import (
	"os"
	"log"

	"gopkg.in/yaml.v2"
)

type OpenAI struct {
	Token string `yaml:"token"`
}

type Azure struct {
	Key          string            `yaml:"key"`
	Endpoint     string            `yaml:"endpoint"`
	ModelMapping map[string]string `yaml:"model-mapping"`
}

type Server struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Config struct {
	Mode   string `yaml:"mode"`
	Server Server `yaml:"server"`
	OpenAI OpenAI `yaml:"openai"`
	Azure  Azure  `yaml:"azure"`
}

func NewConfig(path string) Config {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %s", err)
	}

	return config
}
