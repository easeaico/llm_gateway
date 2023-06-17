package server

import (
	"log"
	"sync/atomic"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type ClientPool struct {
	counter int32
	tn      int
	tokens  []string
	clients []*openai.Client
	limits  []int64
	ticker  *time.Ticker
}

func NewClientPool() *ClientPool {
	p := &ClientPool{
		counter: 0,
		tokens:  config.OpenAI.Tokens,
		tn:      len(config.OpenAI.Tokens),
	}

	var clients []*openai.Client
	for _, token := range p.tokens {
		clients = append(clients, openai.NewClient(token))
	}
	p.clients = clients

	p.limits = make([]int64, p.tn)

	p.ticker = time.NewTicker(3 * time.Second)
	go p.resetLimit()

	return p
}

func (p *ClientPool) GetAvailableClient() *openai.Client {
	index := atomic.AddInt32(&p.counter, 1)
	hash := int(index) % p.tn
	i := hash
	for {
		if p.limits[i] == 0 {
			return p.clients[i]
		}

		i = (i + 1) % p.tn
		if i == hash {
			return p.clients[i]
		}
	}
}

func (p *ClientPool) MarkRateLimit(client *openai.Client) {
	for i, c := range p.clients {
		if c == client {
			atomic.CompareAndSwapInt64(&p.limits[i], p.limits[i], time.Now().Unix())
			return
		}
	}
}

func (p *ClientPool) Close() {
	if p.ticker == nil {
		return
	}

	p.ticker.Stop()
}

func (p *ClientPool) resetLimit() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
		}
	}()

	if p.ticker == nil {
		return
	}

	for {
		select {
		case <-p.ticker.C:
			now := time.Now().Unix()

			for i := 0; i < p.tn; i++ {
				limit := atomic.LoadInt64(&p.limits[i])
				if limit != 0 && now-limit > 60 {
					atomic.CompareAndSwapInt64(&p.limits[i], p.limits[i], 0)
				}
			}
		}
	}
}
