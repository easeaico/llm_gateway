package server_test

import (
	"testing"
	"time"

	"github.com/easeaico/llm_mesh/internal/config"
	"github.com/easeaico/llm_mesh/internal/server"
)

func TestClients(t *testing.T) {
	cfg := &config.Config{
		Mode: "mix",
		OpenAI: config.OpenAI{
			Token: "",
		},
		Azure: config.Azure{
			Key:          "",
			Endpoint:     "",
			ModelMapping: map[string]string{},
		},
		Mix: config.Mix{
			Pipe: []string{"openai", "azure"},
		},
	}

	clients := server.NewClients(cfg)
	client := clients.GetAvailableClient()
	if client == nil {
		t.Errorf("get available client error")
		return
	}

	clients.MarkCurrentRateLimit()
	c := clients.GetAvailableClient()
	if c == nil {
		t.Errorf("get available client error")
		return
	}

	if c == client {
		t.Error("get same client")
		return
	}

	time.Sleep(65 * time.Second)

	c = clients.GetAvailableClient()
	if c != client {
		t.Error("get unknown client")
		return
	}
}

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0

	go func() {
		for {
			select {
			case <-ticker.C:
				count++
			}
		}
	}()

	time.Sleep(5 * time.Second)

	if count != 5 {
		t.Errorf("Expected count to be 5, got %d", count)
	}
}
