# Bugsnag MCP Server (Go)

This project implements a Model Context Protocol (MCP) server in Go for interacting with Bugsnag APIs. It is scaffolded using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features
- Query Bugsnag error and project information via MCP tools (to be implemented)

## Examples
### Get the organizations your user belongs to
```
List the organizations I belong to
```

### Get the projects in an organization
```
list the projects in org "my-org"
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
