package server

import (
	"github.com/easeaico/llm_mesh/internal/config"
	"github.com/easeaico/llm_mesh/pkg/llm_mesh"
	openai "github.com/sashabaranov/go-openai"
)

type chatCompletionServer struct {
	llm_mesh.UnimplementedChatCompletionServiceServer

	mode   Mode
	client *openai.Client
}

func NewChatCompletionServer(conf config.Config) llm_mesh.ChatCompletionServiceServer {
	mode := NewMode(conf)
	client := mode.CreateClient(conf)
	return &chatCompletionServer{
		mode:   mode,
		client: client,
	}
}

func (s *chatCompletionServer) ChatCompletion(req *llm_mesh.ChatCompletionRequest, stream llm_mesh.ChatCompletionService_ChatCompletionServer) error {
	return s.mode.ChatCompletion(s.client, req, stream)
}
