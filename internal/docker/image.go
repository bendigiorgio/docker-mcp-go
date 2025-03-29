package docker

func (c *Client) GetImageTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
func (cli *Client) ImageCreate(ctx context.Context, parentReference string, options image.CreateOptions) (io.ReadCloser, error)
func (cli *Client) ImageHistory(ctx context.Context, imageID string, historyOpts ...ImageHistoryOption) ([]image.HistoryResponseItem, error)
func (cli *Client) ImageImport(ctx context.Context, source image.ImportSource, ref string, ...) (io.ReadCloser, error)
func (cli *Client) ImageInspect(ctx context.Context, imageID string, inspectOpts ...ImageInspectOption) (image.InspectResponse, error)
func (cli *Client) ImageInspectWithRaw(ctx context.Context, imageID string) (image.InspectResponse, []byte, error)deprecated
func (cli *Client) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
func (cli *Client) ImageLoad(ctx context.Context, input io.Reader, loadOpts ...ImageLoadOption) (image.LoadResponse, error)
func (cli *Client) ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)
func (cli *Client) ImagePush(ctx context.Context, image string, options image.PushOptions) (io.ReadCloser, error)
func (cli *Client) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
func (cli *Client) ImageSave(ctx context.Context, imageIDs []string, saveOpts ...ImageSaveOption) (io.ReadCloser, error)
func (cli *Client) ImageSearch(ctx context.Context, term string, options registry.SearchOptions) ([]registry.SearchResult, error)
func (cli *Client) ImageTag(ctx context.Context, source, target string) error
func (cli *Client) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (image.PruneReport, error)
**/
