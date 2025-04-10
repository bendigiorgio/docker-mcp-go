package docker

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initImageTools() {
	// ImageBuild
	imageBuildTool := mcp.NewTool("ImageBuild",
		mcp.WithDescription("Build an image from a Dockerfile"),
		mcp.WithString("Dockerfile", mcp.Required(), mcp.Description("Path to Dockerfile")),
		mcp.WithString("Context", mcp.Required(), mcp.Description("Build context as base64 encoded tar")),
		mcp.WithString("Tags", mcp.Description("Tags to apply to the image")),
		mcp.WithObject("BuildArgs", mcp.Description("Build-time variables")),
		mcp.WithObject("Labels", mcp.Description("Image labels")),
		mcp.WithBoolean("NoCache", mcp.Description("Disable build cache")),
		mcp.WithBoolean("Remove", mcp.Description("Remove intermediate containers")),
		mcp.WithBoolean("ForceRemove", mcp.Description("Always remove intermediate containers")),
	)
	if err := c.RegisterTool(&imageBuildTool, c.ImageBuildHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImageBuild tool")
	}

	// ImagePull
	imagePullTool := mcp.NewTool("ImagePull",
		mcp.WithDescription("Pull an image from a registry"),
		mcp.WithString("Reference", mcp.Required(), mcp.Description("Image reference to pull")),
		mcp.WithString("Platform", mcp.Description("Platform to pull image for")),
		mcp.WithBoolean("All", mcp.Description("Pull all tagged images")),
	)
	if err := c.RegisterTool(&imagePullTool, c.ImagePullHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImagePull tool")
	}

	// ImagePush
	imagePushTool := mcp.NewTool("ImagePush",
		mcp.WithDescription("Push an image to a registry"),
		mcp.WithString("Reference", mcp.Required(), mcp.Description("Image reference to push")),
	)
	if err := c.RegisterTool(&imagePushTool, c.ImagePushHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImagePush tool")
	}

	// ImageList
	imageListTool := mcp.NewTool("ImageList",
		mcp.WithDescription("List images"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing images")),
		mcp.WithBoolean("All", mcp.Description("Show all images")),
	)
	if err := c.RegisterTool(&imageListTool, c.ImageListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImageList tool")
	}

	// ImageRemove
	imageRemoveTool := mcp.NewTool("ImageRemove",
		mcp.WithDescription("Remove an image"),
		mcp.WithString("ImageID", mcp.Required(), mcp.Description("Image ID to remove")),
		mcp.WithBoolean("Force", mcp.Description("Force removal")),
		mcp.WithBoolean("PruneChildren", mcp.Description("Remove untagged parents")),
	)
	if err := c.RegisterTool(&imageRemoveTool, c.ImageRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImageRemove tool")
	}

	// ImageInspect
	imageInspectTool := mcp.NewTool("ImageInspect",
		mcp.WithDescription("Inspect an image"),
		mcp.WithString("ImageID", mcp.Required(), mcp.Description("Image ID to inspect")),
	)
	if err := c.RegisterTool(&imageInspectTool, c.ImageInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImageInspect tool")
	}

	// ImageTag
	imageTagTool := mcp.NewTool("ImageTag",
		mcp.WithDescription("Tag an image"),
		mcp.WithString("Source", mcp.Required(), mcp.Description("Source image ID")),
		mcp.WithString("Target", mcp.Required(), mcp.Description("Target tag name")),
	)
	if err := c.RegisterTool(&imageTagTool, c.ImageTagHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImageTag tool")
	}

	// ImagesPrune
	imagesPruneTool := mcp.NewTool("ImagesPrune",
		mcp.WithDescription("Prune unused images"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when pruning")),
	)
	if err := c.RegisterTool(&imagesPruneTool, c.ImagesPruneHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImagesPrune tool")
	}

	// ImageHistory
	imageHistoryTool := mcp.NewTool("ImageHistory",
		mcp.WithDescription("Show image history"),
		mcp.WithString("ImageID", mcp.Required(), mcp.Description("Image ID to show history for")),
	)
	if err := c.RegisterTool(&imageHistoryTool, c.ImageHistoryHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register ImageHistory tool")
	}
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

type ImageInspectResponse struct {
	Image image.InspectResponse `json:"image"`
	Raw   []byte                `json:"raw,omitempty"`
}

func (c *Client) ImageBuildHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImageBuild")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	// Convert base64 context to tar reader
	contextReader, err := utils.Base64ToTarReader(req.Params.Arguments["Context"].(string))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode build context")
		return nil, err
	}
	defer contextReader.Close()

	options := types.ImageBuildOptions{
		Dockerfile:  req.Params.Arguments["Dockerfile"].(string),
		Tags:        []string{getStringParam(req, "Tags", "")},
		BuildArgs:   convertToPointerStringMap(convertToStringMap(req.Params.Arguments["BuildArgs"])),
		Labels:      convertToStringMap(req.Params.Arguments["Labels"]),
		NoCache:     getBoolParam(req, "NoCache", false),
		Remove:      getBoolParam(req, "Remove", true),
		ForceRemove: getBoolParam(req, "ForceRemove", false),
	}

	response, err := c.cli.ImageBuild(ctx, contextReader, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build image")
		return nil, err
	}
	defer response.Body.Close()

	buildOutput, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read build output")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(buildOutput))
}

func (c *Client) ImagePullHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImagePull")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := image.PullOptions{
		Platform: getStringParam(req, "Platform", ""),
		All:      getBoolParam(req, "All", false),
	}

	ref := req.Params.Arguments["Reference"].(string)
	log.Debug().Str("reference", ref).Msg("Pulling image")

	reader, err := c.cli.ImagePull(ctx, ref, options)
	if err != nil {
		log.Error().Err(err).Str("reference", ref).Msg("Failed to pull image")
		return nil, err
	}
	defer reader.Close()

	pullOutput, err := io.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read pull output")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(pullOutput))
}

func (c *Client) ImagePushHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImagePush")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	ref := req.Params.Arguments["Reference"].(string)
	log.Debug().Str("reference", ref).Msg("Pushing image")

	reader, err := c.cli.ImagePush(ctx, ref, image.PushOptions{})
	if err != nil {
		log.Error().Err(err).Str("reference", ref).Msg("Failed to push image")
		return nil, err
	}
	defer reader.Close()

	pushOutput, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read push output")
		return nil, err
	}

	return utils.FormatSuccessResponse(string(pushOutput))
}

func (c *Client) ImageListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImageList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := image.ListOptions{
		All: getBoolParam(req, "All", false),
	}

	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
					}
				}
			}
		}
	}

	images, err := c.cli.ImageList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list images")
		return nil, err
	}

	return utils.FormatSuccessResponse(images)
}

func (c *Client) ImageRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImageRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := image.RemoveOptions{
		Force:         getBoolParam(req, "Force", false),
		PruneChildren: getBoolParam(req, "PruneChildren", true),
	}

	imageID := req.Params.Arguments["ImageID"].(string)
	log.Debug().Str("imageID", imageID).Msg("Removing image")

	results, err := c.cli.ImageRemove(ctx, imageID, options)
	if err != nil {
		log.Error().Err(err).Str("imageID", imageID).Msg("Failed to remove image")
		return nil, err
	}

	return utils.FormatSuccessResponse(results)
}

func (c *Client) ImageInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImageInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	imageID := req.Params.Arguments["ImageID"].(string)
	log.Debug().Str("imageID", imageID).Msg("Inspecting image")

	img, _, err := c.cli.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		log.Error().Err(err).Str("imageID", imageID).Msg("Failed to inspect image")
		return nil, err
	}

	return utils.FormatSuccessResponse(img)
}

func (c *Client) ImageTagHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImageTag")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	source := req.Params.Arguments["Source"].(string)
	target := req.Params.Arguments["Target"].(string)
	log.Debug().Str("source", source).Str("target", target).Msg("Tagging image")

	if err := c.cli.ImageTag(ctx, source, target); err != nil {
		log.Error().Err(err).Str("source", source).Str("target", target).Msg("Failed to tag image")
		return nil, err
	}

	return utils.FormatSuccessResponse("Image tagged successfully")
}

func (c *Client) ImagesPruneHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImagesPrune")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	pruneFilters := dcf.NewArgs()
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						pruneFilters.Add(key, strValue)
					}
				}
			}
		}
	}

	report, err := c.cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to prune images")
		return nil, err
	}

	return utils.FormatSuccessResponse(report)
}

func (c *Client) ImageHistoryHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("ImageHistory")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	imageID := req.Params.Arguments["ImageID"].(string)
	log.Debug().Str("imageID", imageID).Msg("Getting image history")

	history, err := c.cli.ImageHistory(ctx, imageID)
	if err != nil {
		log.Error().Err(err).Str("imageID", imageID).Msg("Failed to get image history")
		return nil, err
	}

	return utils.FormatSuccessResponse(history)
}

// Helper functions
func convertToPointerStringMap(input map[string]string) map[string]*string {
	result := make(map[string]*string)
	for k, v := range input {
		val := v
		result[k] = &val
	}
	return result
}
