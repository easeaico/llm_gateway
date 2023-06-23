package config

import (
	"io/ioutil"
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

type Mix struct {
	Pipe []string `yaml:"pipe"`
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
	Mix    Mix    `yaml:"mix"`
}

func ReadConfigFile(path string) Config {
	var config Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %s", err)
	}

	return config
}
