package docker

func (c *Client) GetPluginTools() []DockerTool {
	return []DockerTool{}
}

/**
func (cli *Client) PluginCreate(ctx context.Context, createContext io.Reader, ...) error
func (cli *Client) PluginDisable(ctx context.Context, name string, options types.PluginDisableOptions) error
func (cli *Client) PluginEnable(ctx context.Context, name string, options types.PluginEnableOptions) error
func (cli *Client) PluginInspectWithRaw(ctx context.Context, name string) (*types.Plugin, []byte, error)
func (cli *Client) PluginInstall(ctx context.Context, name string, options types.PluginInstallOptions) (rc io.ReadCloser, err error)
func (cli *Client) PluginList(ctx context.Context, filter filters.Args) (types.PluginsListResponse, error)
func (cli *Client) PluginPush(ctx context.Context, name string, registryAuth string) (io.ReadCloser, error)
func (cli *Client) PluginRemove(ctx context.Context, name string, options types.PluginRemoveOptions) error
func (cli *Client) PluginSet(ctx context.Context, name string, args []string) error
func (cli *Client) PluginUpgrade(ctx context.Context, name string, options types.PluginInstallOptions) (io.ReadCloser, error)
**/
