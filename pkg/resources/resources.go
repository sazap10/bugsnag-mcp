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

func NewOrganizationResource() mcp.Resource {
	return mcp.NewResource(
		OrganizationResourceURI,
		"Bugsnag Organizations",
		mcp.WithResourceDescription("Retrieves a list of Bugsnag organizations"),
		mcp.WithMIMEType("application/json"),
	)
}

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

func NewProjectResource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		ProjectTemplateURI,
		"Bugsnag Project",
		mcp.WithTemplateDescription("Retrieves a Bugsnag project by ID"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

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

func NewEventResource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		EventTemplateURI,
		"Bugsnag Event",
		mcp.WithTemplateDescription("Retrieves a Bugsnag event by ID"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

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

// extractIDsFromURI Extracts IDs from a URI given a list of segment names (e.g., "projects", "events")
func extractIDsFromURI(uri string, segments ...string) (map[string]string, error) {
	result := make(map[string]string)
	parts := strings.Split(uri, "/")
	for i, part := range parts {
		for _, seg := range segments {
			if part == seg && i+1 < len(parts) {
				result[seg] = parts[i+1]
			}
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no IDs found in URI: %s", uri)
	}
	return result, nil
}
