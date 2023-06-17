package server

import (
	"errors"
	"io"
	"log"

	"github.com/easeaico/llm_mesh/pkg/llm_mesh"
	openai "github.com/sashabaranov/go-openai"
)

type chatCompletionServer struct {
	llm_mesh.UnimplementedChatCompletionServiceServer
}

var (
	clientPool = NewClientPool()
)

func (s chatCompletionServer) ChatCompletion(req *llm_mesh.ChatCompletionRequest, stream llm_mesh.ChatCompletionService_ChatCompletionServer) error {
	var messages []openai.ChatCompletionMessage
	for _, msg := range req.Messages {
		m := openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
			Name:    msg.Name,
		}
		messages = append(messages, m)
	}

	logitBias := make(map[string]int)
	for k, v := range req.LogitBias {
		logitBias[k] = int(v)
	}

	oReq := openai.ChatCompletionRequest{
		Model:            req.Model,
		Messages:         messages,
		MaxTokens:        int(req.MaxTokens),
		Temperature:      req.Temperature,
		TopP:             req.TopP,
		N:                int(req.N),
		Stream:           req.Stream,
		Stop:             req.Stop,
		PresencePenalty:  req.PresencePenalty,
		FrequencyPenalty: req.FrequencyPenalty,
		LogitBias:        logitBias,
		User:             req.User,
	}

	ctx := stream.Context()
	client := clientPool.GetAvailableClient()
	oStream, err := client.CreateChatCompletionStream(ctx, oReq)
	if err != nil {
		log.Printf("ChatCompletionStream error: %v\n", err)
		return err
	}

	for {
		resp, err := oStream.Recv()
		if err == io.EOF {
			return err
		}

		e := &openai.APIError{}
		if errors.As(err, &e) && e.HTTPStatusCode == 429 {
			clientPool.MarkRateLimit(client)
			log.Printf("Stream error: %v", err)
			return err
		}

		if err != nil {
			log.Printf("Stream error: %v", err)
			return err
		}

		var choices []*llm_mesh.ChatCompletionStreamChoice
		for _, choice := range resp.Choices {
			c := &llm_mesh.ChatCompletionStreamChoice{
				Index: int64(choice.Index),
				Delta: &llm_mesh.ChatCompletionStreamChoiceDelta{
					Content: choice.Delta.Content,
					Role:    choice.Delta.Role,
				},
				FinishReason: choice.FinishReason,
			}
			choices = append(choices, c)
		}

		reply := &llm_mesh.ChatCompletionResponse{
			Response: &llm_mesh.ChatCompletionResponse_Stream{
				Stream: &llm_mesh.ChatCompletionStreamResponse{
					Id:      resp.ID,
					Object:  resp.Object,
					Created: resp.Created,
					Model:   resp.Model,
					Choices: choices,
				},
			},
		}

		if err := stream.Send(reply); err != nil {
			log.Printf("Stream error: %v", err)
			return err
		}
	}
}

func NewServer() llm_mesh.ChatCompletionServiceServer {
	return &chatCompletionServer{}
}
