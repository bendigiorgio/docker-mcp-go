package docker

func (c *Client) GetContainerTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) ContainerAttach(ctx context.Context, containerID string, options container.AttachOptions) (types.HijackedResponse, error)
func (cli *Client) ContainerCommit(ctx context.Context, containerID string, options container.CommitOptions) (container.CommitResponse, error)
func (cli *Client) ContainerCreate(ctx context.Context, config *container.Config, ...) (container.CreateResponse, error)
func (cli *Client) ContainerDiff(ctx context.Context, containerID string) ([]container.FilesystemChange, error)
func (cli *Client) ContainerExecAttach(ctx context.Context, execID string, config container.ExecAttachOptions) (types.HijackedResponse, error)
func (cli *Client) ContainerExecCreate(ctx context.Context, containerID string, options container.ExecOptions) (container.ExecCreateResponse, error)
func (cli *Client) ContainerExecInspect(ctx context.Context, execID string) (container.ExecInspect, error)
func (cli *Client) ContainerExecResize(ctx context.Context, execID string, options container.ResizeOptions) error
func (cli *Client) ContainerExecStart(ctx context.Context, execID string, config container.ExecStartOptions) error
func (cli *Client) ContainerExport(ctx context.Context, containerID string) (io.ReadCloser, error)
func (cli *Client) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error)
func (cli *Client) ContainerInspectWithRaw(ctx context.Context, containerID string, getSize bool) (container.InspectResponse, []byte, error)
func (cli *Client) ContainerKill(ctx context.Context, containerID, signal string) error
func (cli *Client) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
func (cli *Client) ContainerLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error)
func (cli *Client) ContainerPause(ctx context.Context, containerID string) error
func (cli *Client) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
func (cli *Client) ContainerRename(ctx context.Context, containerID, newContainerName string) error
func (cli *Client) ContainerResize(ctx context.Context, containerID string, options container.ResizeOptions) error
func (cli *Client) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
func (cli *Client) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
func (cli *Client) ContainerStatPath(ctx context.Context, containerID, path string) (container.PathStat, error)
func (cli *Client) ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponseReader, error)
func (cli *Client) ContainerStatsOneShot(ctx context.Context, containerID string) (container.StatsResponseReader, error)
func (cli *Client) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
func (cli *Client) ContainerTop(ctx context.Context, containerID string, arguments []string) (container.TopResponse, error)
func (cli *Client) ContainerUnpause(ctx context.Context, containerID string) error
func (cli *Client) ContainerUpdate(ctx context.Context, containerID string, updateConfig container.UpdateConfig) (container.UpdateResponse, error)
func (cli *Client) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error)
func (cli *Client) ContainersPrune(ctx context.Context, pruneFilters filters.Args) (container.PruneReport, error)
func (cli *Client) CopyFromContainer(ctx context.Context, containerID, srcPath string) (io.ReadCloser, container.PathStat, error)
func (cli *Client) CopyToContainer(ctx context.Context, containerID, dstPath string, content io.Reader, ...) error
*/
