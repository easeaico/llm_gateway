package server

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/easeaico/llm_mesh/internal/config"
	"github.com/sashabaranov/go-openai"
)

type Clients struct {
	clients []*openai.Client
	mode    string
	total   int64
	current int64
	status  []int64
}

func NewClients(conf *config.Config) *Clients {
	var clients []*openai.Client
	if conf.Mode == "mix" {
		for _, mode := range conf.Mix.Pipe {
			client := createClientByMode(mode, conf)
			clients = append(clients, client)
		}
	} else {
		client := createClientByMode(conf.Mode, conf)
		clients = append(clients, client)
	}

	total := len(clients)
	status := make([]int64, total)
	c := &Clients{
		mode:    conf.Mode,
		clients: clients,
		total:   int64(total),
		current: 0,
		status:  status,
	}
	go c.resetAvailable()

	return c
}

func createClientByMode(mode string, conf *config.Config) *openai.Client {
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
		log.Panic(fmt.Sprintf("unknown mode: %s", mode))
	}

	return nil
}

func (c *Clients) resetAvailable() {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			for i := range c.status {
				st := atomic.LoadInt64(&c.status[i])
				// rate limtting of minute
				if st > 0 && now-st >= 60 {
					atomic.StoreInt64(&c.status[i], 0)
				}
			}

			// reset index
			for i := range c.status {
				st := atomic.LoadInt64(&c.status[i])
				if st == 0 {
					atomic.StoreInt64(&c.current, int64(i))
					break
				}
			}
		}
	}
}

func (c *Clients) MarkCurrentRateLimit() {
	n := atomic.LoadInt64(&c.current)
	i := n % c.total
	ts := time.Now().Unix()
	atomic.StoreInt64(&c.status[i], ts)
	cur := atomic.AddInt64(&c.current, 1)
	log.Printf("current client index is %d", cur)
}

func (c *Clients) GetAvailableClient() *openai.Client {
	switch c.mode {
	case "openai":
		fallthrough
	case "azure":
		return c.clients[0]
	case "mix":
		n := atomic.LoadInt64(&c.current)
		i := n % c.total
		return c.clients[i]
	default:
		log.Panic(fmt.Sprintf("unknown mode: %s", c.mode))
	}

	return nil
}
