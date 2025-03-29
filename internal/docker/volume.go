package docker

import (
	"context"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initVolumeTools() {
	// VolumeCreate
	volumeCreateTool := mcp.NewTool("VolumeCreate",
		mcp.WithDescription("Create a new volume"),
		mcp.WithString("Name", mcp.Description("Name of the volume")),
		mcp.WithString("Driver", mcp.Description("Volume driver name")),
		mcp.WithObject("DriverOpts", mcp.Description("Driver-specific options")),
		mcp.WithObject("Labels", mcp.Description("Labels to set on the volume")),
	)
	if err := c.RegisterTool(&volumeCreateTool, c.VolumeCreateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register VolumeCreate tool")
	}

	// VolumeInspect
	volumeInspectTool := mcp.NewTool("VolumeInspect",
		mcp.WithDescription("Inspect a volume"),
		mcp.WithString("VolumeID", mcp.Required(), mcp.Description("ID or name of the volume")),
	)
	if err := c.RegisterTool(&volumeInspectTool, c.VolumeInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register VolumeInspect tool")
	}

	// VolumeList
	volumeListTool := mcp.NewTool("VolumeList",
		mcp.WithDescription("List volumes"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing volumes")),
	)
	if err := c.RegisterTool(&volumeListTool, c.VolumeListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register VolumeList tool")
	}

	// VolumeRemove
	volumeRemoveTool := mcp.NewTool("VolumeRemove",
		mcp.WithDescription("Remove a volume"),
		mcp.WithString("VolumeID", mcp.Required(), mcp.Description("ID or name of the volume to remove")),
		mcp.WithBoolean("Force", mcp.Description("Force removal even if volume is in use")),
	)
	if err := c.RegisterTool(&volumeRemoveTool, c.VolumeRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register VolumeRemove tool")
	}

	// VolumesPrune
	volumesPruneTool := mcp.NewTool("VolumesPrune",
		mcp.WithDescription("Prune unused volumes"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when pruning volumes")),
	)
	if err := c.RegisterTool(&volumesPruneTool, c.VolumesPruneHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register VolumesPrune tool")
	}
}

/** Volume API methods
func (cli *Client) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error)
func (cli *Client) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error)
func (cli *Client) VolumeInspectWithRaw(ctx context.Context, volumeID string) (volume.Volume, []byte, error)
func (cli *Client) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
func (cli *Client) VolumeRemove(ctx context.Context, volumeID string, force bool) error
func (cli *Client) VolumeUpdate(ctx context.Context, volumeID string, version swarm.Version, ...) error
func (cli *Client) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error)
**/

type VolumeInspectResponse struct {
	Volume volume.Volume `json:"volume"`
	Raw    []byte        `json:"raw,omitempty"`
}

func (c *Client) VolumeCreateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("VolumeCreate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := volume.CreateOptions{
		Name:       getStringParam(req, "Name", ""),
		Driver:     getStringParam(req, "Driver", "local"),
		DriverOpts: convertToStringMap(req.Params.Arguments["DriverOpts"]),
		Labels:     convertToStringMap(req.Params.Arguments["Labels"]),
	}

	vol, err := c.cli.VolumeCreate(ctx, options)
	if err != nil {
		log.Error().Err(err).Str("volume", options.Name).Msg("Failed to create volume")
		return nil, err
	}

	return utils.FormatSuccessResponse(vol)
}

func (c *Client) VolumeInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("VolumeInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	volumeID := req.Params.Arguments["VolumeID"].(string)
	log.Debug().Str("volumeID", volumeID).Msg("Inspecting volume")

	vol, _, err := c.cli.VolumeInspectWithRaw(ctx, volumeID)
	if err != nil {
		log.Error().Err(err).Str("volumeID", volumeID).Msg("Failed to inspect volume")
		return nil, err
	}

	return utils.FormatSuccessResponse(vol)
}

func (c *Client) VolumeListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("VolumeList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := volume.ListOptions{}
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
						log.Debug().Str("filter", key).Str("value", strValue).Msg("Adding volume filter")
					}
				}
			}
		}
	}

	volumes, err := c.cli.VolumeList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list volumes")
		return nil, err
	}

	return utils.FormatSuccessResponse(volumes)
}

func (c *Client) VolumeRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("VolumeRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	volumeID := req.Params.Arguments["VolumeID"].(string)
	force := getBoolParam(req, "Force", false)
	log.Debug().Str("volumeID", volumeID).Bool("force", force).Msg("Removing volume")

	if err := c.cli.VolumeRemove(ctx, volumeID, force); err != nil {
		log.Error().Err(err).Str("volumeID", volumeID).Msg("Failed to remove volume")
		return nil, err
	}

	return utils.FormatSuccessResponse("Volume removed successfully")
}

func (c *Client) VolumesPruneHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("VolumesPrune")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	pruneFilters := dcf.NewArgs()
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						pruneFilters.Add(key, strValue)
					}
				}
			}
		}
	}

	report, err := c.cli.VolumesPrune(ctx, pruneFilters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prune volumes")
		return nil, err
	}

	return utils.FormatSuccessResponse(report)
}

func convertToStringMap(input interface{}) map[string]string {
	result := make(map[string]string)
	if inputMap, ok := input.(map[string]interface{}); ok {
		for k, v := range inputMap {
			if strVal, ok := v.(string); ok {
				result[k] = strVal
			}
		}
	}
	return result
}
