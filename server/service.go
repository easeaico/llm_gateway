package server

import (
	"errors"
	"io"
	"log"

	"github.com/easeaico/llm_mesh/pkg/llm_mesh"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chatCompletionServer struct {
	llm_mesh.UnimplementedChatCompletionServiceServer
	clients *Clients
}

func NewChatCompletionServer(conf *Config) llm_mesh.ChatCompletionServiceServer {
	return &chatCompletionServer{
		clients: NewClients(conf),
	}
}

func (s *chatCompletionServer) ChatCompletion(req *llm_mesh.ChatCompletionRequest, stream llm_mesh.ChatCompletionService_ChatCompletionServer) error {
	var messages []openai.ChatCompletionMessage
	for _, msg := range req.Messages {
		m := openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
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
	client := s.clients.GetAvailableClient()
	oStream, err := client.CreateChatCompletionStream(ctx, oReq)
	if err != nil {
		log.Printf("ChatCompletionStream error: %v\n", err)
		return status.Errorf(codes.Internal, "ChatCompletionStream error: %v", err)
	}

	for {
		resp, err := oStream.Recv()
		if err == io.EOF {
			return nil
		}

		e := &openai.APIError{}
		if errors.As(err, &e) {
			switch e.HTTPStatusCode {
			case 401:
				// invalid auth or key (do not retry)
				log.Printf("Invalid auth error: %v", err)
				return status.Errorf(codes.Unauthenticated, "Invalid auth error: %v", err)
			case 429:
				// rate limiting or engine overload (wait and retry)
				log.Printf("Stream error: %v", err)
				return status.Errorf(codes.ResourceExhausted, "Stream error: %v", err)
			case 500:
				// openai server error (retry)
				log.Printf("Server auth error: %v", err)
				return status.Errorf(codes.Internal, "Invalid auth error: %v", err)
			default:
				// unhandled
			}
		}

		if err != nil {
			log.Printf("Stream error: %v", err)
			return status.Errorf(codes.Internal, "Stream error: %v", err)
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

		reply := &llm_mesh.ChatCompletionStreamResponse{
			Id:      resp.ID,
			Object:  resp.Object,
			Created: resp.Created,
			Model:   resp.Model,
			Choices: choices,
		}
		if err := stream.Send(reply); err != nil {
			log.Printf("Stream error: %v", err)
			return status.Errorf(codes.Internal, "Stream error: %v", err)
		}
	}
}
