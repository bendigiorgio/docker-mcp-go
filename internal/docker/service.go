package docker

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

/**
func (cli *Client) ServerVersion(ctx context.Context) (types.Version, error)
func (cli *Client) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, ...) (swarm.ServiceCreateResponse, error)
func (cli *Client) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error)
func (cli *Client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error)
func (cli *Client) ServiceLogs(ctx context.Context, serviceID string, options container.LogsOptions) (io.ReadCloser, error)
func (cli *Client) ServiceRemove(ctx context.Context, serviceID string) error
func (cli *Client) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, ...) (swarm.ServiceUpdateResponse, error)
**/

func (c *Client) initServiceTools() {
	// ServiceCreate
	serviceCreateTool := mcp.NewTool("ServiceCreate",
		mcp.WithDescription("Create a new service"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the service")),
		mcp.WithObject("TaskTemplate", mcp.Required(), mcp.Description("Task template specification")),
		mcp.WithObject("Mode", mcp.Description("Service mode (replicated or global)")),
		mcp.WithObject("UpdateConfig", mcp.Description("Update configuration")),
		mcp.WithObject("RollbackConfig", mcp.Description("Rollback configuration")),
		mcp.WithObject("Networks", mcp.Description("Network attachments")),
		mcp.WithObject("EndpointSpec", mcp.Description("Endpoint specification")),
		mcp.WithObject("Labels", mcp.Description("Service labels")),
	)
	if err := c.RegisterTool(&serviceCreateTool, c.ServiceCreateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ServiceCreate tool")
	}

	// ServiceInspect
	serviceInspectTool := mcp.NewTool("ServiceInspect",
		mcp.WithDescription("Inspect a service"),
		mcp.WithString("ServiceID", mcp.Required(), mcp.Description("ID of the service to inspect")),
		mcp.WithBoolean("InsertDefaults", mcp.Description("Include default values in response")),
	)
	if err := c.RegisterTool(&serviceInspectTool, c.ServiceInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ServiceInspect tool")
	}

	// ServiceList
	serviceListTool := mcp.NewTool("ServiceList",
		mcp.WithDescription("List services"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing services")),
		mcp.WithBoolean("Status", mcp.Description("Include service status")),
	)
	if err := c.RegisterTool(&serviceListTool, c.ServiceListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ServiceList tool")
	}

	// ServiceLogs
	serviceLogsTool := mcp.NewTool("ServiceLogs",
		mcp.WithDescription("Get service logs"),
		mcp.WithString("ServiceID", mcp.Required(), mcp.Description("ID of the service")),
		mcp.WithBoolean("Stdout", mcp.Description("Show stdout logs")),
		mcp.WithBoolean("Stderr", mcp.Description("Show stderr logs")),
		mcp.WithBoolean("Timestamps", mcp.Description("Show timestamps")),
		mcp.WithBoolean("Follow", mcp.Description("Follow log output")),
		mcp.WithBoolean("Details", mcp.Description("Show extra details")),
	)
	if err := c.RegisterTool(&serviceLogsTool, c.ServiceLogsHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ServiceLogs tool")
	}

	// ServiceRemove
	serviceRemoveTool := mcp.NewTool("ServiceRemove",
		mcp.WithDescription("Remove a service"),
		mcp.WithString("ServiceID", mcp.Required(), mcp.Description("ID of the service to remove")),
	)
	if err := c.RegisterTool(&serviceRemoveTool, c.ServiceRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ServiceRemove tool")
	}

	// ServiceUpdate
	serviceUpdateTool := mcp.NewTool("ServiceUpdate",
		mcp.WithDescription("Update a service"),
		mcp.WithString("ServiceID", mcp.Required(), mcp.Description("ID of the service to update")),
		mcp.WithString("Version", mcp.Required(), mcp.Description("Current version of the service")),
		mcp.WithObject("Spec", mcp.Required(), mcp.Description("New service specification")),
		mcp.WithObject("Rollback", mcp.Description("Rollback specification")),
	)
	if err := c.RegisterTool(&serviceUpdateTool, c.ServiceUpdateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ServiceUpdate tool")
	}
}

type ServiceInspectResponse struct {
	Service swarm.Service `json:"service"`
	Raw     []byte        `json:"raw,omitempty"`
}

func (c *Client) ServiceCreateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ServiceCreate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	var spec swarm.ServiceSpec
	if err := convertMapToStruct(req.Params.Arguments, &spec); err != nil {
		log.Error().Err(err).Msg("Failed to convert service spec")
		return nil, err
	}

	response, err := c.cli.ServiceCreate(ctx, spec, types.ServiceCreateOptions{})
	if err != nil {
		log.Error().Err(err).Str("service", spec.Name).Msg("Failed to create service")
		return nil, err
	}

	return utils.FormatSuccessResponse(response)
}

func (c *Client) ServiceInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ServiceInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	serviceID := req.Params.Arguments["ServiceID"].(string)
	options := types.ServiceInspectOptions{
		InsertDefaults: getBoolParam(req, "InsertDefaults", false),
	}

	service, raw, err := c.cli.ServiceInspectWithRaw(ctx, serviceID, options)
	if err != nil {
		log.Error().Err(err).Str("serviceID", serviceID).Msg("Failed to inspect service")
		return nil, err
	}

	return utils.FormatSuccessResponse(ServiceInspectResponse{
		Service: service,
		Raw:     raw,
	})
}

func (c *Client) ServiceListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ServiceList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := types.ServiceListOptions{
		Status: getBoolParam(req, "Status", false),
	}

	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
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

	services, err := c.cli.ServiceList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list services")
		return nil, err
	}

	return utils.FormatSuccessResponse(services)
}

func (c *Client) ServiceLogsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ServiceLogs")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	serviceID := req.Params.Arguments["ServiceID"].(string)
	options := container.LogsOptions{
		ShowStdout: getBoolParam(req, "Stdout", true),
		ShowStderr: getBoolParam(req, "Stderr", true),
		Timestamps: getBoolParam(req, "Timestamps", false),
		Follow:     getBoolParam(req, "Follow", false),
		Details:    getBoolParam(req, "Details", false),
	}

	reader, err := c.cli.ServiceLogs(ctx, serviceID, options)
	if err != nil {
		log.Error().Err(err).Str("serviceID", serviceID).Msg("Failed to get service logs")
		return nil, err
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Str("serviceID", serviceID).Msg("Failed to read service logs")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(logs))
}

func (c *Client) ServiceRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ServiceRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	serviceID := req.Params.Arguments["ServiceID"].(string)
	log.Debug().Str("serviceID", serviceID).Msg("Removing service")

	if err := c.cli.ServiceRemove(ctx, serviceID); err != nil {
		log.Error().Err(err).Str("serviceID", serviceID).Msg("Failed to remove service")
		return nil, err
	}

	return utils.FormatSuccessResponse("Service removed successfully")
}

func (c *Client) ServiceUpdateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ServiceUpdate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	serviceID := req.Params.Arguments["ServiceID"].(string)
	var version swarm.Version
	versionIndex, err := strconv.ParseUint(req.Params.Arguments["Version"].(string), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse service version")
		return nil, err
	}
	version.Index = versionIndex

	var spec swarm.ServiceSpec
	if err := convertMapToStruct(req.Params.Arguments["Spec"].(map[string]interface{}), &spec); err != nil {
		log.Error().Err(err).Msg("Failed to convert service spec")
		return nil, err
	}

	response, err := c.cli.ServiceUpdate(ctx, serviceID, version, spec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Error().Err(err).Str("serviceID", serviceID).Msg("Failed to update service")
		return nil, err
	}

	return utils.FormatSuccessResponse(response)
}

// Helper functions
func convertMapToStruct(input map[string]interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}
