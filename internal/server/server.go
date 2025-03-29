package server

import (
	"github.com/bendigiorgio/docker-mcp-go/internal/docker"
	"github.com/mark3labs/mcp-go/server"
	mc "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

func InitializeMCPServer(name string, version string) (*mc.MCPServer, error) {
	s := mc.NewMCPServer(
		name,
		version,
		mc.WithResourceCapabilities(true, true),
		mc.WithLogging(),
		server.WithPromptCapabilities(true),
	)
	dcClient, err := docker.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Docker client")
		return nil, err
	}

	tools := dcClient.GetAllTools()
	for _, tool := range tools {
		s.AddTool(tool.Tool, tool.Handler)
	}

	return s, nil

}
