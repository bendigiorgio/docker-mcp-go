package docker

func (c *Client) GetNodeTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) NodeInspectWithRaw(ctx context.Context, nodeID string) (swarm.Node, []byte, error)
func (cli *Client) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error)
func (cli *Client) NodeRemove(ctx context.Context, nodeID string, options types.NodeRemoveOptions) error
func (cli *Client) NodeUpdate(ctx context.Context, nodeID string, version swarm.Version, node swarm.NodeSpec) error
**/
