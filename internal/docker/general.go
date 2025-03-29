package docker

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) GetGeneralTools() []DockerTool {
	return []DockerTool{
		{
			Tool:    mcp.NewTool("ClientVersion", mcp.WithDescription("Returns the client version")),
			Handler: c.ClientVersionHandler,
		},
		{
			Tool:    mcp.NewTool("Ping", mcp.WithDescription("Ping the Docker server")),
			Handler: c.PingHandler,
		},
	}
}

/** Overall tools:
func (cli *Client) BuildCachePrune(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error)
func (cli *Client) BuildCancel(ctx context.Context, id string) error
func (cli *Client) CheckpointCreate(ctx context.Context, containerID string, options checkpoint.CreateOptions) error
func (cli *Client) CheckpointDelete(ctx context.Context, containerID string, options checkpoint.DeleteOptions) error
func (cli *Client) CheckpointList(ctx context.Context, container string, options checkpoint.ListOptions) ([]checkpoint.Summary, error)
func (cli *Client) ClientVersion() string
func (cli *Client) Close() error
func (cli *Client) DaemonHost() string
func (cli *Client) DialHijack(ctx context.Context, url, proto string, meta map[string][]string) (net.Conn, error)
func (cli *Client) Dialer() func(context.Context) (net.Conn, error)
func (cli *Client) DiskUsage(ctx context.Context, options types.DiskUsageOptions) (types.DiskUsage, error)
func (cli *Client) DistributionInspect(ctx context.Context, imageRef, encodedRegistryAuth string) (registry.DistributionInspect, error)
func (cli *Client) Events(ctx context.Context, options events.ListOptions) (<-chan events.Message, <-chan error)
func (cli *Client) HTTPClient() *http.Client
func (cli *Client) Info(ctx context.Context) (system.Info, error)
func (cli *Client) NegotiateAPIVersion(ctx context.Context)
func (cli *Client) NegotiateAPIVersionPing(pingResponse types.Ping)
func (cli *Client) NewVersionError(ctx context.Context, APIrequired, feature string) error
func (cli *Client) Ping(ctx context.Context) (types.Ping, error)
func (cli *Client) RegistryLogin(ctx context.Context, auth registry.AuthConfig) (registry.AuthenticateOKBody, error)
**/

func (c *Client) ClientVersionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res := c.cli.ClientVersion()
	return mcp.NewToolResultText(res), nil
}

func (c *Client) PingHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res, err := c.cli.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error pinging Docker")
		return nil, err
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling Ping response")
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonRes)), nil
}
