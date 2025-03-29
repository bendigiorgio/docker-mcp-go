package docker

func (c *Client) GetVolumeTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error)
func (cli *Client) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error)
func (cli *Client) VolumeInspectWithRaw(ctx context.Context, volumeID string) (volume.Volume, []byte, error)
func (cli *Client) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
func (cli *Client) VolumeRemove(ctx context.Context, volumeID string, force bool) error
func (cli *Client) VolumeUpdate(ctx context.Context, volumeID string, version swarm.Version, ...) error
func (cli *Client) VolumesPrune(ctx context.Context, pruneFilters filters.Args) (volume.PruneReport, error)
**/
