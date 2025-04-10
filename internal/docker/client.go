package docker

import (
	"context"
	"errors"
	"sync"

	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type Client struct {
	cli *client.Client

	// Tool management
	tools     map[string]DockerTool
	toolsLock sync.RWMutex
}

type DockerTool struct {
	Definition *mcp.Tool
	Handler    ToolHandler
}

type ToolHandler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	c := &Client{
		cli:   cli,
		tools: make(map[string]DockerTool),
	}

	c.initAllTools()
	return c, nil
}

func (c *Client) initAllTools() {
	c.initConfigTools()
	c.initTaskTools()
	c.initSwarmTools()
	c.initSecretTools()
	c.initNodeTools()
	c.initVolumeTools()
	// c.initContainerTools()
	c.initPluginTools()
	c.initServiceTools()
	c.initNetworkTools()
	c.initImageTools()
	// c.initGeneralTools()
}

// RegisterTool is now type-safe and validates input
func (c *Client) RegisterTool(tool *mcp.Tool, handler ToolHandler) error {
	if tool == nil {
		return ErrNilTool
	}
	if handler == nil {
		return ErrNilHandler
	}

	c.toolsLock.Lock()
	defer c.toolsLock.Unlock()

	if _, exists := c.tools[tool.Name]; exists {
		return ErrToolExists
	}

	c.tools[tool.Name] = DockerTool{
		Definition: tool,
		Handler:    handler,
	}
	return nil
}

// GetTool returns both the definition and handler
func (c *Client) GetTool(name string) (DockerTool, bool) {
	c.toolsLock.RLock()
	defer c.toolsLock.RUnlock()

	tool, exists := c.tools[name]
	return tool, exists
}

// GetAllTools returns a copy of all registered tools
func (c *Client) GetAllTools() map[string]DockerTool {
	c.toolsLock.RLock()
	defer c.toolsLock.RUnlock()

	// Return a copy to prevent external modification
	toolsCopy := make(map[string]DockerTool, len(c.tools))
	for k, v := range c.tools {
		toolsCopy[k] = v
	}
	return toolsCopy
}

// Custom errors
var (
	ErrNilTool    = errors.New("tool definition cannot be nil")
	ErrNilHandler = errors.New("tool handler cannot be nil")
	ErrToolExists = errors.New("tool already exists")
)
