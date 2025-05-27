package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sazap10/bugsnag-mcp/pkg/config"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	OrganizationResourceURI = "bugsnag://organizations"
	ProjectTemplateURI      = "bugsnag://projects/{id}"
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
