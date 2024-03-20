package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type OpenAI struct {
	Token string `yaml:"token"`
}

type Azure struct {
	ModelMapping map[string]string `yaml:"model-mapping"`
	Key          string            `yaml:"key"`
	Endpoint     string            `yaml:"endpoint"`
}

type Server struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Config struct {
	Azure  Azure  `yaml:"azure"`
	Mode   string `yaml:"mode"`
	OpenAI OpenAI `yaml:"openai"`
	Server Server `yaml:"server"`
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
