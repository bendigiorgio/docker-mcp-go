package server

import (
	"fmt"

	"github.com/bendigiorgio/docker-mcp-go/internal/docker"
	mc "github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type DockerMCPServer struct {
	*mc.MCPServer
	DockerClient *docker.Client
}

func InitializeMCPServer(name, version string) (*DockerMCPServer, error) {
	// Initialize base MCP server
	s := mc.NewMCPServer(
		name,
		version,
		mc.WithResourceCapabilities(true, true),
		mc.WithLogging(),
		mc.WithPromptCapabilities(true),
	)

	// Initialize Docker client
	dcClient, err := docker.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Register all Docker tools
	if err := registerDockerTools(s, dcClient); err != nil {
		return nil, fmt.Errorf("failed to register Docker tools: %w", err)
	}

	return &DockerMCPServer{
		MCPServer:    s,
		DockerClient: dcClient,
	}, nil
}

func registerDockerTools(s *mc.MCPServer, dc *docker.Client) error {
	tools := dc.GetAllTools()
	for name, tool := range tools {
		s.AddTool(*tool.Definition, mc.ToolHandlerFunc(tool.Handler))
		log.Debug().
			Str("tool", name).
			Msg("Successfully registered Docker tool")
	}
	return nil
}
