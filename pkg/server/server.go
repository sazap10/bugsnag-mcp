package server

import (
	"context"
	"log/slog"
	"os"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/sazap10/bugsnag-mcp/pkg/config"
	"github.com/sazap10/bugsnag-mcp/pkg/resources"
	"github.com/sazap10/bugsnag-mcp/pkg/tools"
)

// NewMCPServer creates a new MCP server with the given name, version, and configuration.
func NewMCPServer(name, version string, cfg *config.Config, hooks ...*mcpserver.Hooks) *mcpserver.MCPServer {
	opts := []mcpserver.ServerOption{
		mcpserver.WithResourceCapabilities(true, true),
		mcpserver.WithToolCapabilities(true),
		mcpserver.WithLogging(),
	}

	// Add hooks if provided
	for _, hook := range hooks {
		opts = append(opts, mcpserver.WithHooks(hook))
	}

	// Create the MCP server
	server := mcpserver.NewMCPServer(name, version, opts...)

	// Register the resources
	registerResources(server, cfg)

	// Register the tools
	registerTools(server, cfg)

	return server
}

// registerResources registers the resources with the MCP server.
func registerResources(server *mcpserver.MCPServer, cfg *config.Config) {
	// Add the organization resource
	orgResource := resources.NewOrganizationResource()
	server.AddResource(orgResource, resources.HandleOrganizationResource(cfg))
	// Add the project resource template
	projectResource := resources.NewProjectResource()
	server.AddResourceTemplate(projectResource, resources.HandleProjectResource(cfg))
	// Add the event resource template
	eventResource := resources.NewEventResource()
	server.AddResourceTemplate(eventResource, resources.HandleEventResource(cfg))
}

// registerTools registers the tools with the MCP server.
func registerTools(server *mcpserver.MCPServer, cfg *config.Config) {
	orgTool := tools.NewGetUserOrganizationsTool()
	server.AddTool(orgTool, tools.HandleGetUserOrganizationsTool(cfg))

	projectTool := tools.NewGetUserProjectsTool()
	server.AddTool(projectTool, tools.HandleGetUserProjectsTool(cfg))

	eventTool := tools.NewGetProjectEventTool()
	server.AddTool(eventTool, tools.HandleGetProjectEventTool(cfg))

	eventsTool := tools.NewGetProjectEventsTool()
	server.AddTool(eventsTool, tools.HandleGetProjectEventsTool(cfg))
}

// ServeStdio starts the MCP server with stdio transport.
func ServeStdio(ctx context.Context, server *mcpserver.MCPServer) error {
	// Create a new stdio transport
	stdioTransport := mcpserver.NewStdioServer(server)

	// create context function
	// This function can be used to set up the context for the server
	// It can be used to set up authentication, logging, etc.
	contextFunc := func(ctx context.Context) context.Context {
		return ctx
	}
	// Set the context function for the server
	stdioTransport.SetContextFunc(contextFunc)

	// Start the server with the stdio transport
	return stdioTransport.Listen(ctx, os.Stdin, os.Stdout)
}

// ServeSSE starts the MCP server with SSE transport.
func ServeSSE(ctx context.Context, server *mcpserver.MCPServer, addr string) error {
	sseServer := mcpserver.NewSSEServer(server)

	//start the server with the SSE transport
	slog.Info("Starting SSE server", slog.String("address", addr))
	return sseServer.Start(addr)
}
