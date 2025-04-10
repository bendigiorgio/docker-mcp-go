package docker

import (
	"context"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initNetworkTools() {
	// NetworkCreate
	networkCreateTool := mcp.NewTool("NetworkCreate",
		mcp.WithDescription("Create a new network"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the network")),
		mcp.WithString("Driver", mcp.Description("Network driver to use")),
		mcp.WithObject("Options", mcp.Description("Driver-specific options")),
		mcp.WithObject("IPAM", mcp.Description("IPAM configuration")),
		mcp.WithBoolean("Internal", mcp.Description("Restrict external access")),
		mcp.WithBoolean("Attachable", mcp.Description("Enable manual container attachment")),
		mcp.WithBoolean("Ingress", mcp.Description("Create swarm routing-mesh network")),
		mcp.WithObject("Labels", mcp.Description("Labels to set on the network")),
	)
	if err := c.RegisterTool(&networkCreateTool, c.NetworkCreateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworkCreate tool")
	}

	// NetworkConnect
	networkConnectTool := mcp.NewTool("NetworkConnect",
		mcp.WithDescription("Connect a container to a network"),
		mcp.WithString("NetworkID", mcp.Required(), mcp.Description("Network ID or name")),
		mcp.WithString("ContainerID", mcp.Required(), mcp.Description("Container ID or name")),
		mcp.WithObject("EndpointSettings", mcp.Description("Endpoint configuration")),
	)
	if err := c.RegisterTool(&networkConnectTool, c.NetworkConnectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworkConnect tool")
	}

	// NetworkDisconnect
	networkDisconnectTool := mcp.NewTool("NetworkDisconnect",
		mcp.WithDescription("Disconnect a container from a network"),
		mcp.WithString("NetworkID", mcp.Required(), mcp.Description("Network ID or name")),
		mcp.WithString("ContainerID", mcp.Required(), mcp.Description("Container ID or name")),
		mcp.WithBoolean("Force", mcp.Description("Force disconnection")),
	)
	if err := c.RegisterTool(&networkDisconnectTool, c.NetworkDisconnectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworkDisconnect tool")
	}

	// NetworkInspect
	networkInspectTool := mcp.NewTool("NetworkInspect",
		mcp.WithDescription("Inspect a network"),
		mcp.WithString("NetworkID", mcp.Required(), mcp.Description("Network ID or name")),
		mcp.WithBoolean("Verbose", mcp.Description("Detailed inspect output")),
		mcp.WithBoolean("Scope", mcp.Description("Filter by scope (swarm, local, etc)")),
	)
	if err := c.RegisterTool(&networkInspectTool, c.NetworkInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworkInspect tool")
	}

	// NetworkList
	networkListTool := mcp.NewTool("NetworkList",
		mcp.WithDescription("List networks"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing networks")),
	)
	if err := c.RegisterTool(&networkListTool, c.NetworkListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworkList tool")
	}

	// NetworkRemove
	networkRemoveTool := mcp.NewTool("NetworkRemove",
		mcp.WithDescription("Remove a network"),
		mcp.WithString("NetworkID", mcp.Required(), mcp.Description("Network ID or name")),
	)
	if err := c.RegisterTool(&networkRemoveTool, c.NetworkRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworkRemove tool")
	}

	// NetworksPrune
	networksPruneTool := mcp.NewTool("NetworksPrune",
		mcp.WithDescription("Prune unused networks"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when pruning networks")),
	)
	if err := c.RegisterTool(&networksPruneTool, c.NetworksPruneHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NetworksPrune tool")
	}
}

/** Network API methods
func (cli *Client) NetworkConnect(ctx context.Context, networkID, containerID string, ...) error
func (cli *Client) NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
func (cli *Client) NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error
func (cli *Client) NetworkInspect(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, error)
func (cli *Client) NetworkInspectWithRaw(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, []byte, error)
func (cli *Client) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
func (cli *Client) NetworkRemove(ctx context.Context, networkID string) error
func (cli *Client) NetworksPrune(ctx context.Context, pruneFilters filters.Args) (network.PruneReport, error)
**/

type NetworkInspectResponse struct {
	Network network.Inspect `json:"network"`
	Raw     []byte          `json:"raw,omitempty"`
}

func (c *Client) NetworkCreateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworkCreate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	name := req.Params.Arguments["Name"].(string)

	options := network.CreateOptions{
		Driver:     getStringParam(req, "Driver", "bridge"),
		Options:    convertToStringMap(req.Params.Arguments["Options"]),
		Internal:   getBoolParam(req, "Internal", false),
		Attachable: getBoolParam(req, "Attachable", false),
		Ingress:    getBoolParam(req, "Ingress", false),
		Labels:     convertToStringMap(req.Params.Arguments["Labels"]),
	}

	if ipam, ok := req.Params.Arguments["IPAM"].(map[string]interface{}); ok {
		options.IPAM = convertToIPAM(ipam)
	}

	res, err := c.cli.NetworkCreate(ctx, name, options)
	if err != nil {
		log.Error().Err(err).Str("network", name).Msg("Failed to create network")
		return nil, err
	}

	return utils.FormatSuccessResponse(res)
}

func (c *Client) NetworkConnectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworkConnect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	networkID := req.Params.Arguments["NetworkID"].(string)
	containerID := req.Params.Arguments["ContainerID"].(string)

	var endpointSettings *network.EndpointSettings
	if settings, ok := req.Params.Arguments["EndpointSettings"].(map[string]interface{}); ok {
		endpointSettings = convertToEndpointSettings(settings)
	}

	if err := c.cli.NetworkConnect(ctx, networkID, containerID, endpointSettings); err != nil {
		log.Error().Err(err).
			Str("networkID", networkID).
			Str("containerID", containerID).
			Msg("Failed to connect container to network")
		return nil, err
	}

	return utils.FormatSuccessResponse("Container connected to network successfully")
}

func (c *Client) NetworkDisconnectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworkDisconnect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	networkID := req.Params.Arguments["NetworkID"].(string)
	containerID := req.Params.Arguments["ContainerID"].(string)
	force := getBoolParam(req, "Force", false)

	if err := c.cli.NetworkDisconnect(ctx, networkID, containerID, force); err != nil {
		log.Error().Err(err).
			Str("networkID", networkID).
			Str("containerID", containerID).
			Bool("force", force).
			Msg("Failed to disconnect container from network")
		return nil, err
	}

	return utils.FormatSuccessResponse("Container disconnected from network successfully")
}

func (c *Client) NetworkInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworkInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	networkID := req.Params.Arguments["NetworkID"].(string)
	options := network.InspectOptions{
		Verbose: getBoolParam(req, "Verbose", false),
		Scope:   getStringParam(req, "Scope", ""),
	}

	net, _, err := c.cli.NetworkInspectWithRaw(ctx, networkID, options)
	if err != nil {
		log.Error().Err(err).Str("networkID", networkID).Msg("Failed to inspect network")
		return nil, err
	}

	return utils.FormatSuccessResponse(net)
}

func (c *Client) NetworkListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworkList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := network.ListOptions{}
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
						log.Debug().Str("filter", key).Str("value", strValue).Msg("Adding network filter")
					}
				}
			}
		}
	}

	networks, err := c.cli.NetworkList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list networks")
		return nil, err
	}

	return utils.FormatSuccessResponse(networks)
}

func (c *Client) NetworkRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworkRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	networkID := req.Params.Arguments["NetworkID"].(string)
	log.Debug().Str("networkID", networkID).Msg("Removing network")

	if err := c.cli.NetworkRemove(ctx, networkID); err != nil {
		log.Error().Err(err).Str("networkID", networkID).Msg("Failed to remove network")
		return nil, err
	}

	return utils.FormatSuccessResponse("Network removed successfully")
}

func (c *Client) NetworksPruneHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NetworksPrune")
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

	report, err := c.cli.NetworksPrune(ctx, pruneFilters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prune networks")
		return nil, err
	}

	return utils.FormatSuccessResponse(report)
}

// Helper functions for network specific conversions
func convertToIPAM(ipam map[string]interface{}) *network.IPAM {
	result := &network.IPAM{
		Driver:  getStringFromMap(ipam, "Driver", "default"),
		Options: convertToStringMap(ipam["Options"]),
	}

	if configs, ok := ipam["Config"].([]interface{}); ok {
		for _, config := range configs {
			if configMap, ok := config.(map[string]interface{}); ok {
				ipamConfig := network.IPAMConfig{
					Subnet:     getStringFromMap(configMap, "Subnet", ""),
					IPRange:    getStringFromMap(configMap, "IPRange", ""),
					Gateway:    getStringFromMap(configMap, "Gateway", ""),
					AuxAddress: convertToStringMap(configMap["AuxAddress"]),
				}
				result.Config = append(result.Config, ipamConfig)
			}
		}
	}

	return result
}

func convertToEndpointSettings(settings map[string]interface{}) *network.EndpointSettings {
	return &network.EndpointSettings{
		IPAMConfig: &network.EndpointIPAMConfig{
			IPv4Address:  getStringFromMap(settings, "IPv4Address", ""),
			IPv6Address:  getStringFromMap(settings, "IPv6Address", ""),
			LinkLocalIPs: convertToStringSlice(settings["LinkLocalIPs"]),
		},
		Links:               convertToStringSlice(settings["Links"]),
		Aliases:             convertToStringSlice(settings["Aliases"]),
		NetworkID:           getStringFromMap(settings, "NetworkID", ""),
		EndpointID:          getStringFromMap(settings, "EndpointID", ""),
		Gateway:             getStringFromMap(settings, "Gateway", ""),
		IPAddress:           getStringFromMap(settings, "IPAddress", ""),
		IPPrefixLen:         getIntFromMap(settings, "IPPrefixLen", 0),
		IPv6Gateway:         getStringFromMap(settings, "IPv6Gateway", ""),
		GlobalIPv6Address:   getStringFromMap(settings, "GlobalIPv6Address", ""),
		GlobalIPv6PrefixLen: getIntFromMap(settings, "GlobalIPv6PrefixLen", 0),
		MacAddress:          getStringFromMap(settings, "MacAddress", ""),
		DriverOpts:          convertToStringMap(settings["DriverOpts"]),
	}
}

func getStringFromMap(m map[string]interface{}, key string, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

func getIntFromMap(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key].(int); ok {
		return val
	}
	return defaultValue
}

func convertToStringSlice(input interface{}) []string {
	var result []string
	if slice, ok := input.([]interface{}); ok {
		for _, item := range slice {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
	}
	return result
}
