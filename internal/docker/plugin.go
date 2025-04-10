package docker

import (
	"context"
	"io"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initPluginTools() {
	// PluginCreate
	pluginCreateTool := mcp.NewTool("PluginCreate",
		mcp.WithDescription("Create a plugin from a tar archive"),
		mcp.WithString("Context", mcp.Required(), mcp.Description("Base64 encoded tar archive of plugin rootfs")),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin")),
	)
	if err := c.RegisterTool(&pluginCreateTool, c.PluginCreateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginCreate tool")
	}

	// PluginEnable
	pluginEnableTool := mcp.NewTool("PluginEnable",
		mcp.WithDescription("Enable a plugin"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin to enable")),
		mcp.WithNumber("Timeout", mcp.Description("Timeout in seconds")),
	)
	if err := c.RegisterTool(&pluginEnableTool, c.PluginEnableHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginEnable tool")
	}

	// PluginDisable
	pluginDisableTool := mcp.NewTool("PluginDisable",
		mcp.WithDescription("Disable a plugin"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin to disable")),
		mcp.WithBoolean("Force", mcp.Description("Force disable even if active")),
	)
	if err := c.RegisterTool(&pluginDisableTool, c.PluginDisableHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginDisable tool")
	}

	// PluginInspect
	pluginInspectTool := mcp.NewTool("PluginInspect",
		mcp.WithDescription("Inspect a plugin"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin to inspect")),
	)
	if err := c.RegisterTool(&pluginInspectTool, c.PluginInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginInspect tool")
	}

	// PluginInstall
	pluginInstallTool := mcp.NewTool("PluginInstall",
		mcp.WithDescription("Install a plugin"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin to install")),
		mcp.WithString("RegistryAuth", mcp.Description("Base64 encoded registry auth config")),
		mcp.WithBoolean("Disabled", mcp.Description("Install but don't enable")),
	)
	if err := c.RegisterTool(&pluginInstallTool, c.PluginInstallHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginInstall tool")
	}

	// PluginList
	pluginListTool := mcp.NewTool("PluginList",
		mcp.WithDescription("List plugins"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing plugins")),
	)
	if err := c.RegisterTool(&pluginListTool, c.PluginListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginList tool")
	}

	// PluginRemove
	pluginRemoveTool := mcp.NewTool("PluginRemove",
		mcp.WithDescription("Remove a plugin"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin to remove")),
		mcp.WithBoolean("Force", mcp.Description("Force removal even if active")),
	)
	if err := c.RegisterTool(&pluginRemoveTool, c.PluginRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginRemove tool")
	}

	// PluginUpgrade
	pluginUpgradeTool := mcp.NewTool("PluginUpgrade",
		mcp.WithDescription("Upgrade a plugin"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the plugin to upgrade")),
		mcp.WithString("RegistryAuth", mcp.Description("Base64 encoded registry auth config")),
		mcp.WithBoolean("Disabled", mcp.Description("Upgrade but don't enable")),
	)
	if err := c.RegisterTool(&pluginUpgradeTool, c.PluginUpgradeHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register PluginUpgrade tool")
	}
}

/**
func (cli *Client) PluginCreate(ctx context.Context, createContext io.Reader, ...) error
func (cli *Client) PluginDisable(ctx context.Context, name string, options types.PluginDisableOptions) error
func (cli *Client) PluginEnable(ctx context.Context, name string, options types.PluginEnableOptions) error
func (cli *Client) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error)
func (cli *Client) PluginInstall(ctx context.Context, name string, options types.PluginInstallOptions) (rc io.ReadCloser, err error)
func (cli *Client) PluginList(ctx context.Context, filter filters.Args) (types.PluginsListResponse, error)
func (cli *Client) PluginPush(ctx context.Context, name string, registryAuth string) (io.ReadCloser, error)
func (cli *Client) PluginRemove(ctx context.Context, name string, options types.PluginRemoveOptions) error
func (cli *Client) PluginSet(ctx context.Context, name string, args []string) error
func (cli *Client) PluginUpgrade(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error)
**/

type PluginInspectResponse struct {
	Plugin *types.Plugin `json:"plugin"`
	Raw    []byte        `json:"raw,omitempty"`
}

func (c *Client) PluginCreateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginCreate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	contextReader, err := utils.Base64ToTarReader(req.Params.Arguments["Context"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode plugin context")
		return nil, err
	}
	defer contextReader.Close()

	name := req.Params.Arguments["Name"].(string)
	if err := c.cli.PluginCreate(ctx, contextReader, types.PluginCreateOptions{RepoName: name}); err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to create plugin")
		return nil, err
	}

	return utils.FormatSuccessResponse("Plugin created successfully")
}

func (c *Client) PluginEnableHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginEnable")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)
	options := types.PluginEnableOptions{
		Timeout: getIntParam(req, "Timeout", 0),
	}

	if err := c.cli.PluginEnable(ctx, name, options); err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to enable plugin")
		return nil, err
	}

	return utils.FormatSuccessResponse("Plugin enabled successfully")
}

func (c *Client) PluginDisableHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginDisable")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)
	options := types.PluginDisableOptions{
		Force: getBoolParam(req, "Force", false),
	}

	if err := c.cli.PluginDisable(ctx, name, options); err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to disable plugin")
		return nil, err
	}

	return utils.FormatSuccessResponse("Plugin disabled successfully")
}

func (c *Client) PluginInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)
	plugin, raw, err := c.cli.PluginInspectWithRaw(ctx, name)
	if err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to inspect plugin")
		return nil, err
	}

	return utils.FormatSuccessResponse(PluginInspectResponse{
		Plugin: plugin,
		Raw:    raw,
	})
}

func (c *Client) PluginInstallHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginInstall")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)
	options := types.PluginInstallOptions{
		Disabled:     getBoolParam(req, "Disabled", false),
		RegistryAuth: getStringParam(req, "RegistryAuth", ""),
	}

	reader, err := c.cli.PluginInstall(ctx, name, options)
	if err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to install plugin")
		return nil, err
	}
	defer reader.Close()

	output, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to read plugin install output")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(output))
}

func (c *Client) PluginListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	filterArgs := filters.NewArgs()
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						filterArgs.Add(key, strValue)
					}
				}
			}
		}
	}

	plugins, err := c.cli.PluginList(ctx, filterArgs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list plugins")
		return nil, err
	}

	return utils.FormatSuccessResponse(plugins)
}

func (c *Client) PluginRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)
	options := types.PluginRemoveOptions{
		Force: getBoolParam(req, "Force", false),
	}

	if err := c.cli.PluginRemove(ctx, name, options); err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to remove plugin")
		return nil, err
	}

	return utils.FormatSuccessResponse("Plugin removed successfully")
}

func (c *Client) PluginUpgradeHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("PluginUpgrade")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)
	options := types.PluginInstallOptions{
		Disabled:     getBoolParam(req, "Disabled", false),
		RegistryAuth: getStringParam(req, "RegistryAuth", ""),
	}

	reader, err := c.cli.PluginUpgrade(ctx, name, options)
	if err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to upgrade plugin")
		return nil, err
	}
	defer reader.Close()

	output, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Str("plugin", name).Msg("Failed to read plugin upgrade output")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(output))
}

// Helper functions
func convertToPluginPrivileges(input interface{}) types.PluginPrivileges {
	var privileges types.PluginPrivileges
	if privs, ok := input.([]interface{}); ok {
		for _, p := range privs {
			if priv, ok := p.(map[string]interface{}); ok {
				privilege := types.PluginPrivilege{
					Name:        getStringFromMap(priv, "Name", ""),
					Description: getStringFromMap(priv, "Description", ""),
					Value:       convertToStringSlice(priv["Value"]),
				}
				privileges = append(privileges, privilege)
			}
		}
	}
	return privileges
}

func getIntParam(req mcp.CallToolRequest, param string, defaultValue int) int {
	if val, ok := req.Params.Arguments[param].(int); ok {
		return val
	}
	return defaultValue
}
