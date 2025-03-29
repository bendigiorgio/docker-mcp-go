package docker

import (
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) GetConfigTools() []DockerTool {
	return []DockerTool{
		{
			Tool: mcp.NewTool("ConfigCreate", mcp.WithDescription("ConfigCreate creates a new config."),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the config")),
				mcp.WithObject("Labels", mcp.Description("Labels of the config (key-value pairs)")),
				mcp.WithString("Data", mcp.Required(), mcp.Description("Data of the config")),
			),
			Handler: c.ConfigCreateHandler,
		},
		{
			Tool: mcp.NewTool("ConfigInspect", mcp.WithDescription("ConfigInspect inspects a config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the config")),
			),
			Handler: c.ConfigInspectHandler,
		},
		{
			Tool: mcp.NewTool("ConfigList", mcp.WithDescription("ConfigList lists configs."),
				mcp.WithObject("Filters", mcp.Description("Filters to list configs (key-value pairs)")),
			),
			Handler: c.ConfigListHandler,
		},
		{
			Tool: mcp.NewTool("ConfigRemove", mcp.WithDescription("ConfigRemove removes a config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the config")),
			),
			Handler: c.ConfigRemoveHandler,
		},
		{
			Tool: mcp.NewTool("ConfigUpdate", mcp.WithDescription("ConfigUpdate updates a config."),
				mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the config")),
				mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the config")),
				mcp.WithObject("Labels", mcp.Description("Labels of the config (key-value pairs)")),
				mcp.WithString("Data", mcp.Required(), mcp.Description("Data of the config")),
			),
			Handler: c.ConfigUpdateHandler,
		},
	}
}

/** Config Tools:
func (cli *Client) ConfigCreate(ctx context.Context, config swarm.ConfigSpec) (types.ConfigCreateResponse, error)
func (cli *Client) ConfigInspectWithRaw(ctx context.Context, id string) (swarm.Config, []byte, error)
func (cli *Client) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error)
func (cli *Client) ConfigRemove(ctx context.Context, id string) error
func (cli *Client) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error
**/

func (c *Client) ConfigCreateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name:   request.Params.Arguments["Name"].(string),
			Labels: request.Params.Arguments["Labels"].(map[string]string),
		},
		Data: []byte(request.Params.Arguments["Data"].(string)),
	}

	res, err := c.cli.ConfigCreate(ctx, config)
	if err != nil {
		log.Error().Err(err).Msg("Error creating config")
		return nil, err
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling config create response")
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonRes)), nil
}

func (c *Client) ConfigInspectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Trace().Msgf("Inspecting config: %s", request.Params.Arguments["ID"].(string))
	res, _, err := c.cli.ConfigInspectWithRaw(ctx, request.Params.Arguments["ID"].(string))
	if err != nil {
		return nil, err
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling config inspect response")
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonRes)), nil
}

func (c *Client) ConfigListHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config := types.ConfigListOptions{}
	if request.Params.Arguments["Filters"] != nil {
		config.Filters = filters.NewArgs()
		for key, values := range request.Params.Arguments["Filters"].(map[string][]interface{}) {
			for _, value := range values {
				log.Debug().Msgf("Adding filter: %s=%s", key, value.(string))
				config.Filters.Add(key, value.(string))
			}
		}
	}
	log.Trace().Msgf("Listing configs: %v", config)
	res, err := c.cli.ConfigList(ctx, config)
	if err != nil {
		log.Error().Err(err).Msg("Error listing configs")
		return nil, err
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling config list response")
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonRes)), nil
}

func (c *Client) ConfigRemoveHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Trace().Msgf("Removing config: %s", request.Params.Arguments["ID"].(string))
	err := c.cli.ConfigRemove(ctx, request.Params.Arguments["ID"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Error removing config")
		return nil, err
	}
	return mcp.NewToolResultText("Config removed"), nil
}

func (c *Client) ConfigUpdateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Trace().Msgf("Updating config: %s", request.Params.Arguments["ID"].(string))
	config := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name:   request.Params.Arguments["Name"].(string),
			Labels: request.Params.Arguments["Labels"].(map[string]string),
		},
		Data: []byte(request.Params.Arguments["Data"].(string)),
	}
	err := c.cli.ConfigUpdate(ctx, request.Params.Arguments["ID"].(string), swarm.Version{}, config)
	if err != nil {
		log.Error().Err(err).Msg("Error updating config")
		return nil, err
	}
	return mcp.NewToolResultText("Config updated"), nil
}
