package docker

import (
	"context"
	"strconv"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initSwarmTools() {
	// SwarmInit
	swarmInitTool := mcp.NewTool("SwarmInit",
		mcp.WithDescription("Initialize a new swarm cluster"),
		mcp.WithString("ListenAddr", mcp.Required(), mcp.Description("Listen address for swarm manager")),
		mcp.WithString("AdvertiseAddr", mcp.Description("Advertise address for swarm manager")),
		mcp.WithBoolean("ForceNewCluster", mcp.Description("Force creation of a new cluster")),
		mcp.WithObject("Spec", mcp.Description("Swarm specification")),
	)
	if err := c.RegisterTool(&swarmInitTool, c.SwarmInitHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmInit tool")
	}

	// SwarmJoin
	swarmJoinTool := mcp.NewTool("SwarmJoin",
		mcp.WithDescription("Join a node to an existing swarm"),
		mcp.WithString("ListenAddr", mcp.Required(), mcp.Description("Listen address for swarm node")),
		mcp.WithString("AdvertiseAddr", mcp.Description("Advertise address for swarm node")),
		mcp.WithString("RemoteAddr", mcp.Required(), mcp.Description("Address of manager node to join")),
		mcp.WithString("JoinToken", mcp.Required(), mcp.Description("Secret token for joining")),
	)
	if err := c.RegisterTool(&swarmJoinTool, c.SwarmJoinHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmJoin tool")
	}

	// SwarmLeave
	swarmLeaveTool := mcp.NewTool("SwarmLeave",
		mcp.WithDescription("Leave the swarm"),
		mcp.WithBoolean("Force", mcp.Description("Force leave even if this is the last manager")),
	)
	if err := c.RegisterTool(&swarmLeaveTool, c.SwarmLeaveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmLeave tool")
	}

	// SwarmInspect
	swarmInspectTool := mcp.NewTool("SwarmInspect",
		mcp.WithDescription("Inspect the current swarm cluster"),
	)
	if err := c.RegisterTool(&swarmInspectTool, c.SwarmInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmInspect tool")
	}

	// SwarmUpdate
	swarmUpdateTool := mcp.NewTool("SwarmUpdate",
		mcp.WithDescription("Update swarm cluster configuration"),
		mcp.WithObject("Spec", mcp.Required(), mcp.Description("New swarm specification")),
		mcp.WithString("Version", mcp.Required(), mcp.Description("Current version of the swarm")),
	)
	if err := c.RegisterTool(&swarmUpdateTool, c.SwarmUpdateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmUpdate tool")
	}

	// SwarmGetUnlockKey
	swarmUnlockKeyTool := mcp.NewTool("SwarmGetUnlockKey",
		mcp.WithDescription("Get the swarm unlock key"),
	)
	if err := c.RegisterTool(&swarmUnlockKeyTool, c.SwarmGetUnlockKeyHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmGetUnlockKey tool")
	}

	// SwarmUnlock
	swarmUnlockTool := mcp.NewTool("SwarmUnlock",
		mcp.WithDescription("Unlock a locked swarm"),
		mcp.WithString("UnlockKey", mcp.Required(), mcp.Description("Secret unlock key")),
	)
	if err := c.RegisterTool(&swarmUnlockTool, c.SwarmUnlockHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SwarmUnlock tool")
	}
}

/** Swarm API methods
func (cli *Client) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error)
func (cli *Client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error)
func (cli *Client) SwarmInspect(ctx context.Context) (swarm.Swarm, error)
func (cli *Client) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error
func (cli *Client) SwarmLeave(ctx context.Context, force bool) error
func (cli *Client) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error
func (cli *Client) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, ...) error
**/

func (c *Client) SwarmInitHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmInit")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	initReq := swarm.InitRequest{
		ListenAddr:      req.Params.Arguments["ListenAddr"].(string),
		AdvertiseAddr:   getStringParam(req, "AdvertiseAddr", ""),
		ForceNewCluster: getBoolParam(req, "ForceNewCluster", false),
	}

	if spec, ok := req.Params.Arguments["Spec"].(map[string]interface{}); ok {
		initReq.Spec = convertMapToSwarmSpec(spec)
	}

	nodeID, err := c.cli.SwarmInit(ctx, initReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize swarm")
		return nil, err
	}

	return utils.FormatSuccessResponse(map[string]string{"NodeID": nodeID})
}

func (c *Client) SwarmJoinHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmJoin")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	joinReq := swarm.JoinRequest{
		ListenAddr:    req.Params.Arguments["ListenAddr"].(string),
		AdvertiseAddr: getStringParam(req, "AdvertiseAddr", ""),
		RemoteAddrs:   []string{req.Params.Arguments["RemoteAddr"].(string)},
		JoinToken:     req.Params.Arguments["JoinToken"].(string),
	}

	if err := c.cli.SwarmJoin(ctx, joinReq); err != nil {
		log.Error().Err(err).Msg("Failed to join swarm")
		return nil, err
	}

	return utils.FormatSuccessResponse("Successfully joined swarm")
}

func (c *Client) SwarmLeaveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmLeave")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	force := getBoolParam(req, "Force", false)
	if err := c.cli.SwarmLeave(ctx, force); err != nil {
		log.Error().Err(err).Bool("force", force).Msg("Failed to leave swarm")
		return nil, err
	}

	return utils.FormatSuccessResponse("Successfully left swarm")
}

func (c *Client) SwarmInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	swarmInfo, err := c.cli.SwarmInspect(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to inspect swarm")
		return nil, err
	}

	return utils.FormatSuccessResponse(swarmInfo)
}

func (c *Client) SwarmUpdateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmUpdate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	var version swarm.Version
	versionIndex, err := strconv.ParseUint(req.Params.Arguments["Version"].(string), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse swarm version")
		return nil, err
	}
	version.Index = versionIndex

	var spec swarm.Spec
	if specMap, ok := req.Params.Arguments["Spec"].(map[string]interface{}); ok {
		spec = convertMapToSwarmSpec(specMap)
	}

	// TODO update flags
	if err := c.cli.SwarmUpdate(ctx, version, spec, swarm.UpdateFlags{}); err != nil {
		log.Error().Err(err).Msg("Failed to update swarm")
		return nil, err
	}

	return utils.FormatSuccessResponse("Swarm updated successfully")
}

func (c *Client) SwarmGetUnlockKeyHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmGetUnlockKey")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	unlockKey, err := c.cli.SwarmGetUnlockKey(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get swarm unlock key")
		return nil, err
	}

	return utils.FormatSuccessResponse(unlockKey)
}

func (c *Client) SwarmUnlockHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SwarmUnlock")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	unlockReq := swarm.UnlockRequest{
		UnlockKey: req.Params.Arguments["UnlockKey"].(string),
	}

	if err := c.cli.SwarmUnlock(ctx, unlockReq); err != nil {
		log.Error().Err(err).Msg("Failed to unlock swarm")
		return nil, err
	}

	return utils.FormatSuccessResponse("Swarm unlocked successfully")
}

// Helper functions
func convertMapToSwarmSpec(specMap map[string]interface{}) swarm.Spec {
	// Implement conversion logic from map to swarm.Spec
	// This is a simplified version - you may need to expand it based on your needs
	var spec swarm.Spec
	if val, ok := specMap["Name"].(string); ok {
		spec.Name = val
	}
	// Add other fields as needed
	return spec
}

func getStringParam(req mcp.CallToolRequest, param string, defaultValue string) string {
	if val, ok := req.Params.Arguments[param].(string); ok {
		return val
	}
	return defaultValue
}
