package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/sazap10/bugsnag-mcp/pkg/config"
)

// Tool IDs
const (
	GetUserOrganizationsToolID = "get_user_organizations"
	GetUserProjectsToolID      = "get_user_projects"
)

func NewGetUserOrganizationsTool() mcp.Tool {
	return mcp.NewTool(
		GetUserOrganizationsToolID,
		mcp.WithDescription("Retrieves the organizations for the current user from Bugsnag"),
	)
}

func HandleGetUserOrganizationsTool(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		orgs, _, err := cfg.APIClient.CurrentUser.ListOrganizations(ctx, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to retrieve organizations: %v", err)), nil
		}

		orgsJSON, err := json.MarshalIndent(orgs, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal organizations: %v", err)), nil
		}

		return mcp.NewToolResultText(string(orgsJSON)), nil
	}
}

func NewGetUserProjectsTool() mcp.Tool {
	return mcp.NewTool(
		GetUserProjectsToolID,
		mcp.WithDescription("Retrieves the projects for the current user from Bugsnag"),
		mcp.WithString(
			"organization_id",
			mcp.Required(),
			mcp.Description("The ID of the organization to retrieve projects for"),
		),
	)
}

func HandleGetUserProjectsTool(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		org_id, err := req.RequireString("organization_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'organization_id': %v", err)), nil
		}
		projects, _, err := cfg.APIClient.CurrentUser.ListProjects(ctx, org_id, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to retrieve projects: %v", err)), nil
		}

		projectsJSON, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal projects: %v", err)), nil
		}

		return mcp.NewToolResultText(string(projectsJSON)), nil
	}
}
