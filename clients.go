package main

import (
	"log"

	"github.com/sashabaranov/go-openai"
)

type Clients struct {
	mode    string
	clients []*openai.Client
	status  []int64
	total   int64
	current int64
}

func NewClients(conf *Config) *Clients {
	var clients []*openai.Client
	client := createClientByMode(conf.Mode, conf)
	clients = append(clients, client)

	total := len(clients)
	status := make([]int64, total)
	c := &Clients{
		mode:    conf.Mode,
		clients: clients,
		total:   int64(total),
		current: 0,
		status:  status,
	}

	return c
}

func createClientByMode(mode string, conf *Config) *openai.Client {
	switch mode {
	case "openai":
		client := openai.NewClient(conf.OpenAI.Token)
		return client
	case "azure":
		cfg := openai.DefaultAzureConfig(conf.Azure.Key, conf.Azure.Endpoint)
		cfg.AzureModelMapperFunc = func(model string) string {
			return conf.Azure.ModelMapping[model]
		}
		client := openai.NewClientWithConfig(cfg)
		return client
	default:
		log.Panicf("unknown mode: %s", mode)
	}

	return nil
}

func (c *Clients) GetAvailableClient() *openai.Client {
	switch c.mode {
	case "openai":
		fallthrough
	case "azure":
		return c.clients[0]
	default:
		log.Panicf("unknown mode: %s", c.mode)
	}

	return nil
}
