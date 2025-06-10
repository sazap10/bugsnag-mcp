package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/sazap10/bugsnag-mcp/pkg/config"
)

// Tool IDs
const (
	GetUserOrganizationsToolID = "get_user_organizations"
	GetUserProjectsToolID      = "get_user_projects"
	GetProjectEventToolID      = "get_project_event"
	GetProjectEventsToolID     = "get_project_events"
)

// NewGetUserOrganizationsTool returns the MCP tool for listing Bugsnag organizations for the current user.
func NewGetUserOrganizationsTool() mcp.Tool {
	return mcp.NewTool(
		GetUserOrganizationsToolID,
		mcp.WithDescription("Retrieves the organizations for the current user from Bugsnag"),
	)
}

// HandleGetUserOrganizationsTool handles the tool call to retrieve all organizations for the current user.
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

// NewGetUserProjectsTool returns the MCP tool for listing all projects in a specified organization.
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

// HandleGetUserProjectsTool handles the tool call to retrieve all projects for a given organization.
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

// NewGetProjectEventTool returns the MCP tool for retrieving a specific event for a project.
func NewGetProjectEventTool() mcp.Tool {
	return mcp.NewTool(
		GetProjectEventToolID,
		mcp.WithDescription("Retrieves a specific event for a project from Bugsnag"),
		mcp.WithString(
			"project_id",
			mcp.Required(),
			mcp.Description("The ID of the project to retrieve the event for"),
		),
		mcp.WithString(
			"event_id",
			mcp.Required(),
			mcp.Description("The ID/url of the event to retrieve"),
		),
	)
}

// HandleGetProjectEventTool handles the tool call to retrieve a specific event for a project by event ID or link.
func HandleGetProjectEventTool(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'project_id': %v", err)), nil
		}
		reqParam, err := req.RequireString("event_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'event_id': %v", err)), nil
		}

		eventID, err := getEventIDFromIDOrLink(reqParam)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("invalid event ID or link: %v", err)), nil
		}

		event, _, err := cfg.APIClient.Events.GetEvent(ctx, projectID, eventID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to retrieve event: %v", err)), nil
		}

		eventJSON, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal event: %v", err)), nil
		}

		return mcp.NewToolResultText(string(eventJSON)), nil
	}
}

// NewGetProjectEventsTool returns the MCP tool for listing all events for a project.
func NewGetProjectEventsTool() mcp.Tool {
	return mcp.NewTool(
		GetProjectEventsToolID,
		mcp.WithDescription("Retrieves all events for a project from Bugsnag"),
		mcp.WithString(
			"project_id",
			mcp.Required(),
			mcp.Description("The ID of the project to retrieve events for"),
		),
	)
}

// HandleGetProjectEventsTool handles the tool call to retrieve all events for a given project.
func HandleGetProjectEventsTool(cfg *config.Config) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("missing required parameter 'project_id': %v", err)), nil
		}

		// Fetch events for the project
		events, _, err := cfg.APIClient.Events.ListProjectEvents(ctx, projectID, nil)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to retrieve events: %v", err)), nil
		}

		eventsJSON, err := json.MarshalIndent(events, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal events: %v", err)), nil
		}

		return mcp.NewToolResultText(string(eventsJSON)), nil
	}
}

// getEventIDFromIDOrLink extracts the event ID from either a direct ID or a Bugsnag dashboard link.
// If a link is provided, it parses the URL and returns the event_id query parameter.
func getEventIDFromIDOrLink(idOrLink string) (string, error) {
	// If it's a plain ID (no slashes, no http), just return it
	if !strings.Contains(idOrLink, "/") && !strings.HasPrefix(idOrLink, "http") {
		return idOrLink, nil
	}

	// Otherwise, try to parse as URL and extract event_id param
	u, err := url.Parse(idOrLink)
	if err != nil {
		return "", fmt.Errorf("invalid event link: %w", err)
	}
	q := u.Query()
	eventID := q.Get("event_id")
	if eventID == "" {
		return "", fmt.Errorf("event_id not found in link: %s", idOrLink)
	}
	return eventID, nil
}
