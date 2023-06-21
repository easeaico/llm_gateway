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
	Key      string            `yaml:"key"`
	Endpoint string            `yaml:"endpoint"`
	Models   map[string]string `yaml:"models"`
}

type Config struct {
	Mode   string `yaml:"mode"`
	OpenAI OpenAI `yaml:"openai"`
	Azure  Azure  `yaml:"azure"`
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
