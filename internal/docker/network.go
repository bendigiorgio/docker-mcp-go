package docker

func (c *Client) GetNetworkTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) NetworkConnect(ctx context.Context, networkID, containerID string, ...) error
func (cli *Client) NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
func (cli *Client) NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error
func (cli *Client) NetworkInspect(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, error)
func (cli *Client) NetworkInspectWithRaw(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, []byte, error)
func (cli *Client) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
func (cli *Client) NetworkRemove(ctx context.Context, networkID string) error
func (cli *Client) NetworksPrune(ctx context.Context, pruneFilters filters.Args) (network.PruneReport, error)
**/
