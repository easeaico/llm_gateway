package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type chatCompletionHandler struct {
	clients *Clients
}

func newChatCompletionHandler(conf *Config) *chatCompletionHandler {
	return &chatCompletionHandler{
		clients: NewClients(conf),
	}
}

func (h *chatCompletionHandler) HandleCompletions(c echo.Context) error {
	req := openai.ChatCompletionRequest{}
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	client := h.clients.GetAvailableClient()
	ctx := c.Request().Context()
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}

	r := c.Response()
	r.Header().Set(echo.HeaderContentType, "text/event-stream")
	r.WriteHeader(http.StatusOK)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			resp, err := stream.Recv()
			if err != nil {
				return err
			}

			data, err := json.Marshal(resp)
			if err != nil {
				return err
			}

			fmt.Fprintf(r, "data: %s\n\n", string(data))
			r.Flush()
		}
	}
}
