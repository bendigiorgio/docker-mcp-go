package docker

import (
	"context"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initConfigTools() {
	// ConfigCreate
	configCreateTool := mcp.NewTool("ConfigCreate",
		mcp.WithDescription("Creates a new config"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the config")),
		mcp.WithObject("Labels", mcp.Description("Labels of the config (key-value pairs)")),
		mcp.WithString("Data", mcp.Required(), mcp.Description("Base64-encoded config data")),
	)
	if err := c.RegisterTool(&configCreateTool, c.ConfigCreateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ConfigCreate tool")
	}

	// ConfigInspect
	configInspectTool := mcp.NewTool("ConfigInspect",
		mcp.WithDescription("Inspects an existing config"),
		mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the config to inspect")),
	)
	if err := c.RegisterTool(&configInspectTool, c.ConfigInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ConfigInspect tool")
	}

	// ConfigList
	configListTool := mcp.NewTool("ConfigList",
		mcp.WithDescription("Lists all configs"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing configs")),
	)
	if err := c.RegisterTool(&configListTool, c.ConfigListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ConfigList tool")
	}

	// ConfigRemove
	configRemoveTool := mcp.NewTool("ConfigRemove",
		mcp.WithDescription("Removes a config"),
		mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the config to remove")),
	)
	if err := c.RegisterTool(&configRemoveTool, c.ConfigRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ConfigRemove tool")
	}

	// ConfigUpdate
	configUpdateTool := mcp.NewTool("ConfigUpdate",
		mcp.WithDescription("Updates an existing config"),
		mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the config to update")),
		mcp.WithString("Name", mcp.Required(), mcp.Description("New name for the config")),
		mcp.WithObject("Labels", mcp.Description("New labels for the config")),
		mcp.WithString("Data", mcp.Required(), mcp.Description("New base64-encoded config data")),
	)
	if err := c.RegisterTool(&configUpdateTool, c.ConfigUpdateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ConfigUpdate tool")
	}
}

/** Config API methods
func (cli *Client) ConfigCreate(ctx context.Context, config swarm.ConfigSpec) (types.ConfigCreateResponse, error)
func (cli *Client) ConfigInspectWithRaw(ctx context.Context, id string) (swarm.Config, []byte, error)
func (cli *Client) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error)
func (cli *Client) ConfigRemove(ctx context.Context, id string) error
func (cli *Client) ConfigUpdate(ctx context.Context, id string, version swarm.Version, config swarm.ConfigSpec) error
**/

func (c *Client) ConfigCreateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get tool definition for validation
	tool, exists := c.GetTool("ConfigCreate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	// Validate parameters
	if err := utils.ValidateRequestParams(request.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	// Create config spec
	config := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name:   request.Params.Arguments["Name"].(string),
			Labels: request.Params.Arguments["Labels"].(map[string]string),
		},
		Data: []byte(request.Params.Arguments["Data"].(string)),
	}

	// Execute Docker operation
	res, err := c.cli.ConfigCreate(ctx, config)
	if err != nil {
		log.Error().Err(err).Str("config", config.Name).Msg("Failed to create config")
		return nil, err
	}

	return utils.FormatSuccessResponse(res)
}

func (c *Client) ConfigInspectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ConfigInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(request.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	configID := request.Params.Arguments["ID"].(string)
	log.Debug().Str("configID", configID).Msg("Inspecting config")

	config, _, err := c.cli.ConfigInspectWithRaw(ctx, configID)
	if err != nil {
		log.Error().Err(err).Str("configID", configID).Msg("Failed to inspect config")
		return nil, err
	}

	return utils.FormatSuccessResponse(config)
}

func (c *Client) ConfigListHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ConfigList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(request.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := types.ConfigListOptions{}
	if filters, ok := request.Params.Arguments["Filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
					}
				}
			}
		}
	}

	configs, err := c.cli.ConfigList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list configs")
		return nil, err
	}

	return utils.FormatSuccessResponse(configs)
}

func (c *Client) ConfigRemoveHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ConfigRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(request.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	configID := request.Params.Arguments["ID"].(string)
	log.Debug().Str("configID", configID).Msg("Removing config")

	if err := c.cli.ConfigRemove(ctx, configID); err != nil {
		log.Error().Err(err).Str("configID", configID).Msg("Failed to remove config")
		return nil, err
	}

	return utils.FormatSuccessResponse("Config removed successfully")
}

func (c *Client) ConfigUpdateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ConfigUpdate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(request.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	configID := request.Params.Arguments["ID"].(string)
	log.Debug().Str("configID", configID).Msg("Updating config")

	// First get current version
	currentConfig, _, err := c.cli.ConfigInspectWithRaw(ctx, configID)
	if err != nil {
		log.Error().Err(err).Str("configID", configID).Msg("Failed to get current config version")
		return nil, err
	}

	// Prepare update
	config := swarm.ConfigSpec{
		Annotations: swarm.Annotations{
			Name:   request.Params.Arguments["Name"].(string),
			Labels: request.Params.Arguments["Labels"].(map[string]string),
		},
		Data: []byte(request.Params.Arguments["Data"].(string)),
	}

	if err := c.cli.ConfigUpdate(ctx, configID, currentConfig.Version, config); err != nil {
		log.Error().Err(err).Str("configID", configID).Msg("Failed to update config")
		return nil, err
	}

	return utils.FormatSuccessResponse("Config updated successfully")
}
