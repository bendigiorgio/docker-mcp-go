package docker

import (
	"context"
	"strconv"

	"github.com/bendigiorgio/docker-mcp-go/internal/utils"
	"github.com/docker/docker/api/types"
	dcf "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

func (c *Client) initSecretTools() {
	// SecretCreate
	secretCreateTool := mcp.NewTool("SecretCreate",
		mcp.WithDescription("Create a new secret"),
		mcp.WithString("Name", mcp.Required(), mcp.Description("Name of the secret")),
		mcp.WithObject("Labels", mcp.Description("Labels for the secret (key-value pairs)")),
		mcp.WithString("Data", mcp.Required(), mcp.Description("Base64-encoded secret data")),
		mcp.WithObject("Driver", mcp.Description("Secret driver configuration")),
		mcp.WithBoolean("Templating", mcp.Description("Enable templating for the secret")),
	)
	if err := c.RegisterTool(&secretCreateTool, c.SecretCreateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SecretCreate tool")
	}

	// SecretInspect
	secretInspectTool := mcp.NewTool("SecretInspect",
		mcp.WithDescription("Inspect a secret"),
		mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the secret to inspect")),
	)
	if err := c.RegisterTool(&secretInspectTool, c.SecretInspectHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SecretInspect tool")
	}

	// SecretList
	secretListTool := mcp.NewTool("SecretList",
		mcp.WithDescription("List secrets"),
		mcp.WithObject("Filters", mcp.Description("Filters to apply when listing secrets")),
	)
	if err := c.RegisterTool(&secretListTool, c.SecretListHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SecretList tool")
	}

	// SecretRemove
	secretRemoveTool := mcp.NewTool("SecretRemove",
		mcp.WithDescription("Remove a secret"),
		mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the secret to remove")),
	)
	if err := c.RegisterTool(&secretRemoveTool, c.SecretRemoveHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SecretRemove tool")
	}

	// SecretUpdate
	secretUpdateTool := mcp.NewTool("SecretUpdate",
		mcp.WithDescription("Update a secret"),
		mcp.WithString("ID", mcp.Required(), mcp.Description("ID of the secret to update")),
		mcp.WithString("Name", mcp.Required(), mcp.Description("New name for the secret")),
		mcp.WithObject("Labels", mcp.Description("New labels for the secret")),
		mcp.WithString("Data", mcp.Required(), mcp.Description("New base64-encoded secret data")),
		mcp.WithString("Version", mcp.Required(), mcp.Description("Current version of the secret")),
	)
	if err := c.RegisterTool(&secretUpdateTool, c.SecretUpdateHandler); err != nil {
		log.Error().Err(err).Msg("Failed to register SecretUpdate tool")
	}
}

/** Secret API methods
func (cli *Client) SecretCreate(ctx context.Context, secret swarm.SecretSpec) (types.SecretCreateResponse, error)
func (cli *Client) SecretInspectWithRaw(ctx context.Context, id string) (swarm.Secret, []byte, error)
func (cli *Client) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error)
func (cli *Client) SecretRemove(ctx context.Context, id string) error
func (cli *Client) SecretUpdate(ctx context.Context, id string, version swarm.Version, secret swarm.SecretSpec) error
**/

type SecretInspectResponse struct {
	Secret swarm.Secret `json:"secret"`
	Raw    []byte       `json:"raw"`
}

func (c *Client) SecretCreateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SecretCreate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name:   req.Params.Arguments["Name"].(string),
			Labels: convertToLabels(req.Params.Arguments["Labels"]),
		},
		Data: []byte(req.Params.Arguments["Data"].(string)),
	}

	if driver, ok := req.Params.Arguments["Driver"].(map[string]interface{}); ok {
		spec.Driver = &swarm.Driver{
			Name:    driver["Name"].(string),
			Options: convertToDriverOptions(driver["Options"]),
		}
	}

	if templating, ok := req.Params.Arguments["Templating"].(bool); ok {
		spec.Templating = &swarm.Driver{
			Name:    "",
			Options: nil,
		}
		if templating {
			spec.Templating.Name = "golang"
		}
	}

	res, err := c.cli.SecretCreate(ctx, spec)
	if err != nil {
		log.Error().Err(err).Str("secret", spec.Name).Msg("Failed to create secret")
		return nil, err
	}

	return utils.FormatSuccessResponse(res)
}

func (c *Client) SecretInspectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SecretInspect")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	secretID := req.Params.Arguments["ID"].(string)
	log.Debug().Str("secretID", secretID).Msg("Inspecting secret")

	secret, raw, err := c.cli.SecretInspectWithRaw(ctx, secretID)
	if err != nil {
		log.Error().Err(err).Str("secretID", secretID).Msg("Failed to inspect secret")
		return nil, err
	}

	return utils.FormatSuccessResponse(SecretInspectResponse{
		Secret: secret,
		Raw:    raw,
	})
}

func (c *Client) SecretListHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SecretList")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	options := types.SecretListOptions{}
	if filters, ok := req.Params.Arguments["Filters"].(map[string]interface{}); ok {
		options.Filters = dcf.NewArgs()
		for key, values := range filters {
			if valuesSlice, ok := values.([]interface{}); ok {
				for _, value := range valuesSlice {
					if strValue, ok := value.(string); ok {
						options.Filters.Add(key, strValue)
						log.Debug().Str("filter", key).Str("value", strValue).Msg("Adding secret filter")
					}
				}
			}
		}
	}

	secrets, err := c.cli.SecretList(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list secrets")
		return nil, err
	}

	return utils.FormatSuccessResponse(secrets)
}

func (c *Client) SecretRemoveHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SecretRemove")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	secretID := req.Params.Arguments["ID"].(string)
	log.Debug().Str("secretID", secretID).Msg("Removing secret")

	if err := c.cli.SecretRemove(ctx, secretID); err != nil {
		log.Error().Err(err).Str("secretID", secretID).Msg("Failed to remove secret")
		return nil, err
	}

	return utils.FormatSuccessResponse("Secret removed successfully")
}

func (c *Client) SecretUpdateHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tool, exists := c.GetTool("SecretUpdate")
	if !exists {
		return nil, utils.ErrToolNotFound
	}

	if err := utils.ValidateRequestParams(req.Params.Arguments, tool.Definition); err != nil {
		return nil, err
	}

	secretID := req.Params.Arguments["ID"].(string)
	log.Debug().Str("secretID", secretID).Msg("Updating secret")

	var version swarm.Version
	versionStr := req.Params.Arguments["Version"].(string)
	versionIndex, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse secret version")
		return nil, err
	}
	version.Index = versionIndex

	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name:   req.Params.Arguments["Name"].(string),
			Labels: convertToLabels(req.Params.Arguments["Labels"]),
		},
		Data: []byte(req.Params.Arguments["Data"].(string)),
	}

	if err := c.cli.SecretUpdate(ctx, secretID, version, spec); err != nil {
		log.Error().Err(err).Str("secretID", secretID).Msg("Failed to update secret")
		return nil, err
	}

	return utils.FormatSuccessResponse("Secret updated successfully")
}

// Helper functions
func convertToLabels(labels interface{}) map[string]string {
	result := make(map[string]string)
	if labelsMap, ok := labels.(map[string]interface{}); ok {
		for k, v := range labelsMap {
			if strVal, ok := v.(string); ok {
				result[k] = strVal
			}
		}
	}
	return result
}

func convertToDriverOptions(options interface{}) map[string]string {
	result := make(map[string]string)
	if optionsMap, ok := options.(map[string]interface{}); ok {
		for k, v := range optionsMap {
			if strVal, ok := v.(string); ok {
				result[k] = strVal
			}
		}
	}
	return result
}
