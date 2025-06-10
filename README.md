# Bugsnag MCP Server (Go)

This project implements a Model Context Protocol (MCP) server in Go for interacting with Bugsnag APIs. It is scaffolded using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features

- Query Bugsnag error and project information via MCP tools
- List organizations, projects, and events from your Bugsnag account
- Retrieve details for specific events and projects

## Tools Available

The following MCP tools are available in this server:

- **GetUserOrganizations**: List the organizations your Bugsnag user belongs to.
- **GetUserProjects**: List all projects in a specified organization. Requires `organization_id`.
- **GetProjectEvents**: List all events for a specified project. Requires `project_id`.
- **GetProjectEvent**: Retrieve details for a specific event in a project. Requires `project_id` and `event_id` (can be an ID or a Bugsnag dashboard link).

## Resources Available

The following MCP resources are available:

- **bugsnag://organizations**: Retrieve all organizations for the current user.
- **bugsnag://projects/{id}**: Retrieve details for a specific project by ID.
- **bugsnag://projects/{project_id}/events/{id}**: Retrieve details for a specific event by project and event ID.

## Examples

### Get the organizations your user belongs to

```
List the organizations I belong to
```

### Get the projects in an organization

```
list the projects in org "my-org"
```

### Get all events for a project

```
list the events for project "my-project" in organization "my-org"
```

### Get details for a specific event

```
get details for event "<EVENT_LINK_FROM_DASHBOARD>" in project "my-project"
```

## Installation

### Prerequisites

- Go 1.24 or later
- BugSnag account with personal access token

### Build from source

1. Clone the repository:
   ```
   git clone https://github.com/sazap10/bugsnag-mcp
   cd bugsnag-mcp
   ```
2. Build the binary:
   ```
   go build -o bugsnag-mcp .
   ```
3. Copy binary to your PATH:
   ```
   cp bugsnag-mcp /usr/local/bin/bugsnag-mcp
   ```

## Usage

### VS Code

Add the following configuration to `.vscode/mcp.json`, depending on the type you want to use:

#### stdio

```
{
  "inputs": [
    {
      "id": "bugsnag_auth_token",
      "type": "promptString",
      "description": "BugSnag Auth Token",
      "password": true
    }
  ],
  "servers": {
    "bugsnag-mcp": {
      "type": "stdio",
      "command": "bugsnag-mcp",
      "args": [],
      "env": {
        "BUGSNAG_AUTH_TOKEN": "${input:bugsnag_auth_token}"
      }
    }
  }
}
```

<!-- #### SSE
```
{
  "servers": {
    "bugsnag-mcp": {
      "type": "sse",
      "url": "http://localhost:8080/sse",
      "env": {
        "BUGSNAG_AUTH_TOKEN": "${input:bugsnag_auth_token}"
      }
    }
  }
}
``` -->

## References

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [mcp-go library](https://github.com/mark3labs/mcp-go)
- [Bugsnag API docs](https://bugsnagapiv2.docs.apiary.io/)
