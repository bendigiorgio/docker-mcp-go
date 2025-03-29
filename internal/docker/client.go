package docker

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type Client struct {
	cli *client.Client
}

type DockerTool struct {
	Tool    mcp.Tool
	Handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) CloseClient() error {
	return c.cli.Close()
}

func (c *Client) GetAllTools() []DockerTool {
	configTools := c.GetConfigTools()

	allTools := append(configTools, c.GetContainerTools()...)
	allTools = append(allTools, c.GetGeneralTools()...)
	allTools = append(allTools, c.GetImageTools()...)
	allTools = append(allTools, c.GetNetworkTools()...)
	allTools = append(allTools, c.GetNodeTools()...)
	allTools = append(allTools, c.GetPluginTools()...)
	allTools = append(allTools, c.GetSecretTools()...)
	allTools = append(allTools, c.GetServiceTools()...)
	allTools = append(allTools, c.GetSwarmTools()...)
	allTools = append(allTools, c.GetTaskTools()...)
	allTools = append(allTools, c.GetVolumeTools()...)
	return allTools
}
