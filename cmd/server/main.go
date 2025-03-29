package main

import (
	"flag"
	"fmt"

	dcServer "github.com/bendigiorgio/docker-mcp-go/internal/server"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

func main() {
	var transport string
	var port string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&port, "p", "8080", "Port to listen on")
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.Parse()

	mcpServer, err := dcServer.InitializeMCPServer("docker-mcp-server", "1.0.0")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize MCP server")
	}
	if transport == "sse" {
		sseServer := server.NewSSEServer(mcpServer.MCPServer, server.WithBaseURL(fmt.Sprintf("http://localhost:%s", port)))
		log.Info().Msgf("Starting SSE server on port %s", port)
		if err := sseServer.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatal().Err(err).Msg("Failed to start MCP server")
		}
	} else {
		log.Info().Msg("Starting STDIO server")
		if err := server.ServeStdio(mcpServer.MCPServer); err != nil {
			log.Fatal().Err(err).Msg("Failed to start MCP server")
		}
	}
}
