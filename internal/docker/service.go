package docker

func (c *Client) GetServiceTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) ServerVersion(ctx context.Context) (types.Version, error)
func (cli *Client) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, ...) (swarm.ServiceCreateResponse, error)
func (cli *Client) ServiceInspectWithRaw(ctx context.Context, serviceID string, opts types.ServiceInspectOptions) (swarm.Service, []byte, error)
func (cli *Client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error)
func (cli *Client) ServiceLogs(ctx context.Context, serviceID string, options container.LogsOptions) (io.ReadCloser, error)
func (cli *Client) ServiceRemove(ctx context.Context, serviceID string) error
func (cli *Client) ServiceUpdate(ctx context.Context, serviceID string, version swarm.Version, ...) (swarm.ServiceUpdateResponse, error)
**/
