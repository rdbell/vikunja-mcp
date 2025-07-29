# Vikunja MCP Server

A Model Context Protocol (MCP) server that provides integration with Vikunja, an open-source to-do application and project management tool.

## Features

- Get tasks from Vikunja projects
- Create new tasks in projects
- Update existing tasks (title, description, completion status)
- Get all projects
- Create new projects

## Setup

### Prerequisites

- Go 1.24.2 or later
- A running Vikunja instance
- A Vikunja API token

### Installation

1. Clone this repository
2. Set up environment variables either by:

   **Option A: Using a .env file (recommended)**

   ```bash
   cp .env-example .env
   # Edit .env with your Vikunja instance details
   ```

   **Option B: Using environment variables**

   ```bash
   export VIKUNJA_URL="https://your-vikunja-instance.com"
   export VIKUNJA_TOKEN="your-api-token"
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Build the server:

   ```bash
   go build -o vikunja-mcp
   ```

### Getting a Vikunja API Token

1. Log into your Vikunja instance
2. Go to Settings > API Tokens
3. Create a new token with appropriate permissions
4. Copy the token for use in the `VIKUNJA_TOKEN` environment variable

## Usage

Run the MCP server:

```bash
./vikunja-mcp
```

The server will start and listen for MCP protocol messages on stdin/stdout.

## Available Tools

### `get_tasks`

Get all tasks from a specific project.

- `project_id` (integer, required): ID of the project

### `create_task`

Create a new task in a project.

- `project_id` (integer, required): ID of the project
- `title` (string, required): Task title
- `description` (string, optional): Task description

### `update_task`

Update an existing task.

- `task_id` (integer, required): ID of the task to update
- `title` (string, optional): New task title
- `description` (string, optional): New task description
- `done` (boolean, optional): Mark task as completed/incomplete

### `get_projects`

Get all projects from Vikunja.
No parameters required.

### `create_project`

Create a new project.

- `title` (string, required): Project title
- `description` (string, optional): Project description

## Integration with MCP Clients

This server can be used with any MCP-compatible client. Configure your client to use this server as an MCP server.

## Project Structure

- `main.go` - Entry point for the MCP server
- `tools.go` - MCP tool definitions and handlers
- `client.go` - Vikunja API client implementation
