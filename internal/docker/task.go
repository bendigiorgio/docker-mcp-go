package docker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) GetTaskTools() []DockerTool {
	return []DockerTool{
		{
			Tool: mcp.NewTool("TaskInspectWithRaw", mcp.WithDescription("TaskInspectWithRaw inspects a task with raw data."),
				mcp.WithString("taskID", mcp.Required(), mcp.Description("ID of the task")),
			),
			Handler: c.TaskInspectWithRawHandler,
		},
		{
			Tool: mcp.NewTool("TaskList", mcp.WithDescription("TaskList lists tasks."),
				mcp.WithObject("filters", mcp.Description("Filters to list tasks (key-value pairs)")),
			),
			Handler: c.TaskListHandler,
		},
		{
			Tool: mcp.NewTool("TaskLogs", mcp.WithDescription("TaskLogs gets logs of a task."),
				mcp.WithString("taskID", mcp.Required(), mcp.Description("ID of the task")),
			),
			Handler: c.TaskLogsHandler,
		},
	}
}

/**
func (cli *Client) TaskInspectWithRaw(ctx context.Context, taskID string) (swarm.Task, []byte, error)
func (cli *Client) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error)
func (cli *Client) TaskLogs(ctx context.Context, taskID string, options container.LogsOptions) (io.ReadCloser, error)
**/

// For JSON marshalling
type TaskInspectWithRawResponse struct {
	Task  swarm.Task `json:"task"`
	Bytes []byte     `json:"bytes"`
}

func (c *Client) TaskInspectWithRawHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Trace().Msgf("TaskInspectWithRawHandler: %v", req.Params.Arguments)
	task, bytes, err := c.cli.TaskInspectWithRaw(ctx, req.Params.Arguments["taskID"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to inspect task")
		return nil, err
	}

	res := TaskInspectWithRawResponse{
		Task:  task,
		Bytes: bytes,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal task inspect response")
		return nil, err
	}

	return mcp.NewToolResultText(string(resBytes)), nil
}

func (c *Client) TaskListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Trace().Msgf("TaskListHandler: %v", req.Params.Arguments)
	options := types.TaskListOptions{
		Filters: func() filters.Args {
			args := filters.NewArgs()
			for key, values := range req.Params.Arguments["filters"].(map[string][]string) {
				for _, value := range values {
					args.Add(key, value)
				}
			}
			return args
		}(),
	}
	tasks, err := c.cli.TaskList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list tasks")
		return nil, err
	}

	resBytes, err := json.Marshal(tasks)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal task list response")
		return nil, err
	}

	return mcp.NewToolResultText(string(resBytes)), nil
}

func (c *Client) TaskLogsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Trace().Msgf("TaskLogsHandler: %v", req.Params.Arguments)
	taskID := req.Params.Arguments["taskID"].(string)
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	}
	reader, err := c.cli.TaskLogs(ctx, taskID, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get task logs")
		return nil, err
	}
	defer reader.Close()

	bytes, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read task logs")
		return nil, err
	}

	return mcp.NewToolResultText(string(bytes)), nil
}
