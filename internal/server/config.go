package server

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var (
	config Config = ReadConfigFile("./conf/conf.yaml")
)

type Config struct {
	OpenAI struct {
		Tokens []string `yaml:"tokens"`
	} `yaml:"openai"`
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
