package utils

import (
	"encoding/json"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog/log"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var (
	ErrToolNotFound = errors.New("requested tool not found")
)

func FormatResponse(data interface{}, err error) (string, error) {
	response := StandardResponse{
		Success: err == nil,
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Data = data
	}

	jsonResponse, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Msg("Failed to marshal response")
		return "", marshalErr
	}

	return string(jsonResponse), nil
}

func FormatSuccessResponse(data interface{}) (*mcp.CallToolResult, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}
