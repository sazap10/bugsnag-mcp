package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sazap10/bugsnag-mcp/pkg/config"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	OrganizationResourceURI = "bugsnag://organizations"
	ProjectTemplateURI      = "bugsnag://projects/{id}"
	EventTemplateURI        = "bugsnag://projects/{project_id}/events/{id}"
)

// NewOrganizationResource returns the MCP resource for listing Bugsnag organizations.
func NewOrganizationResource() mcp.Resource {
	return mcp.NewResource(
		OrganizationResourceURI,
		"Bugsnag Organizations",
		mcp.WithResourceDescription("Retrieves a list of Bugsnag organizations"),
		mcp.WithMIMEType("application/json"),
	)
}

// HandleOrganizationResource handles requests to retrieve all organizations for the current user.
func HandleOrganizationResource(cfg *config.Config) server.ResourceHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Call the Bugsnag API to get the list of organizations
		orgs, _, err := cfg.APIClient.CurrentUser.ListOrganizations(ctx, nil)
		if err != nil {
			return nil, err
		}

		orgsJSON, err := json.Marshal(orgs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal organizations: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      OrganizationResourceURI,
				MIMEType: "application/json",
				Text:     string(orgsJSON),
			},
		}, nil
	}
}

// NewProjectResource returns the MCP resource template for a single Bugsnag project by ID.
func NewProjectResource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		ProjectTemplateURI,
		"Bugsnag Project",
		mcp.WithTemplateDescription("Retrieves a Bugsnag project by ID"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// HandleProjectResource handles requests to retrieve a specific project by ID from Bugsnag.
func HandleProjectResource(cfg *config.Config) server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract the project ID from the request
		uri := req.Params.URI
		if uri == "" {
			return nil, fmt.Errorf("project ID not provided in request")
		}
		ids, err := extractIDsFromURI(uri, "projects")
		if err != nil {
			return nil, fmt.Errorf("failed to extract project ID from URI: %v", err)
		}
		projectID, ok := ids["projects"]
		if !ok {
			return nil, fmt.Errorf("project ID not found in URI: %s", uri)
		}

		// Call the Bugsnag API to get the project details
		project, _, err := cfg.APIClient.Projects.GetProject(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve project: %v", err)
		}

		projectJSON, err := json.Marshal(project)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal project: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(projectJSON),
			},
		}, nil
	}
}

// NewEventResource returns the MCP resource template for a single Bugsnag event by project and event ID.
func NewEventResource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		EventTemplateURI,
		"Bugsnag Event",
		mcp.WithTemplateDescription("Retrieves a Bugsnag event by ID"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

// HandleEventResource handles requests to retrieve a specific event by project and event ID from Bugsnag.
func HandleEventResource(cfg *config.Config) server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		// Extract the event ID from the request
		uri := req.Params.URI
		if uri == "" {
			return nil, fmt.Errorf("event ID not provided in request")
		}
		ids, err := extractIDsFromURI(uri, "projects", "events")
		if err != nil {
			return nil, fmt.Errorf("failed to extract IDs from URI: %v", err)
		}
		projectID, ok := ids["projects"]
		if !ok {
			return nil, fmt.Errorf("project ID not found in URI: %s", uri)
		}
		eventID, ok := ids["events"]
		if !ok {
			return nil, fmt.Errorf("event ID not found in URI: %s", uri)
		}

		// Call the Bugsnag API to get the event details
		event, _, err := cfg.APIClient.Events.GetEvent(ctx, projectID, eventID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve event: %v", err)
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event: %v", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(eventJSON),
			},
		}, nil
	}
}

// extractIDsFromURI extracts IDs from a URI given a list of segment names (e.g., "projects", "events").
// Returns a map of segment name to ID, e.g. {"projects": "123", "events": "456"}.
func extractIDsFromURI(uri string, segments ...string) (map[string]string, error) {
	result := make(map[string]string)
	parts := strings.Split(uri, "/")
	for i, part := range parts {
		for _, seg := range segments {
			if part == seg && i+1 < len(parts) {
				id := parts[i+1]
				if id != "" {
					result[seg] = id
				}
			}
		}
	}
	if len(result) != len(segments) {
		return nil, fmt.Errorf("not all IDs found in URI: %s", uri)
	}
	return result, nil
}
