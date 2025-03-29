package docker

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initNodeTools() {
	// NodeInspect
	nodeInspectTool := mcp.NewTool("NodeInspect",
		mcp.WithDescription("Inspect a swarm node"),
		mcp.WithString("NodeID", mcp.Required(), mcp.Description("ID of the node to inspect")),
	)
	if err := c.RegisterTool(&nodeInspectTool, c.NodeInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NodeInspect tool")
	}

	// NodeList
	nodeListTool := mcp.NewTool("NodeList",
		mcp.WithDescription("List swarm nodes"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing nodes")),
	)
	if err := c.RegisterTool(&nodeListTool, c.NodeListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NodeList tool")
	}

	// NodeRemove
	nodeRemoveTool := mcp.NewTool("NodeRemove",
		mcp.WithDescription("Remove a node from the swarm"),
		mcp.WithString("NodeID", mcp.Required(), mcp.Description("ID of the node to remove")),
		mcp.WithBoolean("Force", mcp.Description("Force remove a node even if it's a manager")),
	)
	if err := c.RegisterTool(&nodeRemoveTool, c.NodeRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NodeRemove tool")
	}

	// NodeUpdate
	nodeUpdateTool := mcp.NewTool("NodeUpdate",
		mcp.WithDescription("Update a swarm node"),
		mcp.WithString("NodeID", mcp.Required(), mcp.Description("ID of the node to update")),
		mcp.WithString("Version", mcp.Required(), mcp.Description("Current version of the node")),
		mcp.WithObject("Spec", mcp.Required(), mcp.Description("New node specification")),
	)
	if err := c.RegisterTool(&nodeUpdateTool, c.NodeUpdateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register NodeUpdate tool")
	}
}

/** Node API methods
func (cli *Client) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error)
func (cli *Client) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error)
func (cli *Client) NodeRemove(ctx context.Context, nodeID string, options types.NodeRemoveOptions) error
func (cli *Client) NodeUpdate(ctx context.Context, nodeID string, version swarm.Version, node swarm.NodeSpec) error
**/

type NodeInspectResponse struct {
	Node swarm.Node `json:"node"`
	Raw  []byte     `json:"raw"`
}

func (c *Client) NodeInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NodeInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	nodeID := req.Params.Arguments["NodeID"].(string)
	log.Debug().Str("nodeID", nodeID).Msg("Inspecting node")

	node, raw, err := c.cli.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		log.Error().Err(err).Str("nodeID", nodeID).Msg("Failed to inspect node")
		return nil, err
	}

	return utils.FormatSuccessResponse(NodeInspectResponse{
		Node: node,
		Raw:  raw,
	})
}

func (c *Client) NodeListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NodeList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := types.NodeListOptions{}
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
						log.Debug().Str("filter", key).Str("value", strValue).Msg("Adding node filter")
					}
				}
			}
		}
	}

	nodes, err := c.cli.NodeList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list nodes")
		return nil, err
	}

	return utils.FormatSuccessResponse(nodes)
}

func (c *Client) NodeRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NodeRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	nodeID := req.Params.Arguments["NodeID"].(string)
	force := getBoolParam(req, "Force", false)
	log.Debug().Str("nodeID", nodeID).Bool("force", force).Msg("Removing node")

	options := types.NodeRemoveOptions{
		Force: force,
	}

	if err := c.cli.NodeRemove(ctx, nodeID, options); err != nil {
		log.Error().Err(err).Str("nodeID", nodeID).Msg("Failed to remove node")
		return nil, err
	}

	return utils.FormatSuccessResponse("Node removed successfully")
}

func (c *Client) NodeUpdateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("NodeUpdate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	nodeID := req.Params.Arguments["NodeID"].(string)
	log.Debug().Str("nodeID", nodeID).Msg("Updating node")

	var version swarm.Version
	versionIndex, err := strconv.ParseUint(req.Params.Arguments["Version"].(string), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse node version")
		return nil, err
	}
	version.Index = versionIndex

	var spec swarm.NodeSpec
	if specMap, ok := req.Params.Arguments["Spec"].(map[string]interface{}); ok {
		specBytes, err := json.Marshal(specMap)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal node spec")
			return nil, err
		}
		if err := json.Unmarshal(specBytes, &spec); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal node spec")
			return nil, err
		}
	}

	if err := c.cli.NodeUpdate(ctx, nodeID, version, spec); err != nil {
		log.Error().Err(err).Str("nodeID", nodeID).Msg("Failed to update node")
		return nil, err
	}

	return utils.FormatSuccessResponse("Node updated successfully")
}
