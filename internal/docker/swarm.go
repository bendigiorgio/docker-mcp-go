package docker

func (c *Client) GetSwarmTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) SwarmGetUnlockKey(ctx context.Context) (types.SwarmUnlockKeyResponse, error)
func (cli *Client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error)
func (cli *Client) SwarmInspect(ctx context.Context) (swarm.Swarm, error)
func (cli *Client) SwarmJoin(ctx context.Context, req swarm.JoinRequest) error
func (cli *Client) SwarmLeave(ctx context.Context, force bool) error
func (cli *Client) SwarmUnlock(ctx context.Context, req swarm.UnlockRequest) error
func (cli *Client) SwarmUpdate(ctx context.Context, version swarm.Version, swarm swarm.Spec, ...) error
**/
