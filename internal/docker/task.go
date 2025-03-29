package docker

import (
	"context"
	"io"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initTaskTools() {
	// TaskInspect
	taskInspectTool := mcp.NewTool("TaskInspect",
		mcp.WithDescription("Inspect a swarm task with raw data"),
		mcp.WithString("taskID", mcp.Required(), mcp.Description("ID of the task to inspect")),
	)
	if err := c.RegisterTool(&taskInspectTool, c.TaskInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register TaskInspect tool")
	}

	// TaskList
	taskListTool := mcp.NewTool("TaskList",
		mcp.WithDescription("List swarm tasks"),
		mcp.WithObject("filters", mcp.Description("Filters to apply when listing tasks")),
	)
	if err := c.RegisterTool(&taskListTool, c.TaskListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register TaskList tool")
	}

	// TaskLogs
	taskLogsTool := mcp.NewTool("TaskLogs",
		mcp.WithDescription("Get logs from a swarm task"),
		mcp.WithString("taskID", mcp.Required(), mcp.Description("ID of the task to get logs from")),
		mcp.WithBoolean("stdout", mcp.Description("Show stdout logs (default: true)")),
		mcp.WithBoolean("stderr", mcp.Description("Show stderr logs (default: true)")),
		mcp.WithBoolean("timestamps", mcp.Description("Show timestamps (default: false)")),
		mcp.WithBoolean("follow", mcp.Description("Follow log output (default: false)")),
	)
	if err := c.RegisterTool(&taskLogsTool, c.TaskLogsHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register TaskLogs tool")
	}
}

/** Task API methods
func (cli *Client) TaskInspectWithRaw(ctx context.Context, taskID string) (swarm.Task, []byte, error)
func (cli *Client) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error)
func (cli *Client) TaskLogs(ctx context.Context, taskID string, options container.LogsOptions) (io.ReadCloser, error)
**/

type TaskInspectResponse struct {
	Task swarm.Task `json:"task"`
	Raw  []byte     `json:"raw"`
}

func (c *Client) TaskInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("TaskInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	taskID := req.Params.Arguments["taskID"].(string)
	log.Debug().Str("taskID", taskID).Msg("Inspecting task")

	task, raw, err := c.cli.TaskInspectWithRaw(ctx, taskID)
	if err != nil {
		log.Error().Err(err).Str("taskID", taskID).Msg("Failed to inspect task")
		return nil, err
	}

	return utils.FormatSuccessResponse(TaskInspectResponse{
		Task: task,
		Raw:  raw,
	})
}

func (c *Client) TaskListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("TaskList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := types.TaskListOptions{}
	if filters, ok := req.Params.Arguments["filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
						log.Debug().Str("filter", key).Str("value", strValue).Msg("Adding task filter")
					}
				}
			}
		}
	}

	tasks, err := c.cli.TaskList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list tasks")
		return nil, err
	}

	return utils.FormatSuccessResponse(tasks)
}

func (c *Client) TaskLogsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("TaskLogs")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	taskID := req.Params.Arguments["taskID"].(string)
	options := container.LogsOptions{
		ShowStdout: getBoolParam(req, "stdout", true),
		ShowStderr: getBoolParam(req, "stderr", true),
		Timestamps: getBoolParam(req, "timestamps", false),
		Follow:     getBoolParam(req, "follow", false),
	}

	log.Debug().
		Str("taskID", taskID).
		Bool("stdout", options.ShowStdout).
		Bool("stderr", options.ShowStderr).
		Bool("timestamps", options.Timestamps).
		Bool("follow", options.Follow).
		Msg("Getting task logs")

	reader, err := c.cli.TaskLogs(ctx, taskID, options)
	if err != nil {
		log.Error().Err(err).Str("taskID", taskID).Msg("Failed to get task logs")
		return nil, err
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Str("taskID", taskID).Msg("Failed to read task logs")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(logs))
}

// Helper function to safely get boolean parameters with default values
func getBoolParam(req mcp.CallToolRequest, param string, defaultValue bool) bool {
	if val, ok := req.Params.Arguments[param].(bool); ok {
		return val
	}
	return defaultValue
}
