package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) {
	getTasksTool := mcp.NewTool("get_todo",
		mcp.WithDescription("Get a list of personal todos"),
	)
	s.AddTool(getTasksTool, getTasksHandler)

	createTaskTool := mcp.NewTool("create_todo",
		mcp.WithDescription("Create a new personal todo"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the task"),
		),
	)
	s.AddTool(createTaskTool, createTaskHandler)
}

func getTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := getVikunjaClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasks, err := client.GetTasks()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get tasks: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Tasks: %+v", tasks)), nil
}

func createTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := getVikunjaClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectID := 1

	title, err := request.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError("title is required and must be a string"), nil
	}

	task, err := client.CreateTask(projectID, title)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create task: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created task: %+v", task)), nil
}
