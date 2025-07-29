package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) {
	getTasksTool := mcp.NewTool("get_tasks",
		mcp.WithDescription("Get tasks from a Vikunja project"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("ID of the project to get tasks from"),
		),
	)
	s.AddTool(getTasksTool, getTasksHandler)

	createTaskTool := mcp.NewTool("create_task",
		mcp.WithDescription("Create a new task in a Vikunja project"),
		mcp.WithNumber("project_id",
			mcp.Required(),
			mcp.Description("ID of the project to create task in"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the task"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the task"),
		),
	)
	s.AddTool(createTaskTool, createTaskHandler)

	updateTaskTool := mcp.NewTool("update_task",
		mcp.WithDescription("Update an existing task in Vikunja"),
		mcp.WithNumber("task_id",
			mcp.Required(),
			mcp.Description("ID of the task to update"),
		),
		mcp.WithString("title",
			mcp.Description("New title of the task"),
		),
		mcp.WithString("description",
			mcp.Description("New description of the task"),
		),
		mcp.WithBoolean("done",
			mcp.Description("Mark task as done or not done"),
		),
	)
	s.AddTool(updateTaskTool, updateTaskHandler)

	getProjectsTool := mcp.NewTool("get_projects",
		mcp.WithDescription("Get all projects from Vikunja"),
	)
	s.AddTool(getProjectsTool, getProjectsHandler)

	createProjectTool := mcp.NewTool("create_project",
		mcp.WithDescription("Create a new project in Vikunja"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("Title of the project"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the project"),
		),
	)
	s.AddTool(createProjectTool, createProjectHandler)
}

func getTasksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := getVikunjaClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projectID, err := request.RequireInt("project_id")
	if err != nil {
		return mcp.NewToolResultError("project_id is required and must be a number"), nil
	}

	tasks, err := client.GetTasks(projectID)
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

	projectID, err := request.RequireInt("project_id")
	if err != nil {
		return mcp.NewToolResultError("project_id is required and must be a number"), nil
	}

	title, err := request.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError("title is required and must be a string"), nil
	}

	description := request.GetString("description", "")

	task, err := client.CreateTask(projectID, title, description)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create task: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created task: %+v", task)), nil
}

func updateTaskHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := getVikunjaClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	taskID, err := request.RequireInt("task_id")
	if err != nil {
		return mcp.NewToolResultError("task_id is required and must be a number"), nil
	}

	updates := make(map[string]interface{})
	if title := request.GetString("title", ""); title != "" {
		updates["title"] = title
	}
	if description := request.GetString("description", ""); description != "" {
		updates["description"] = description
	}

	arguments := request.GetArguments()
	if done, ok := arguments["done"]; ok {
		updates["done"] = done
	}

	task, err := client.UpdateTask(taskID, updates)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to update task: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Updated task: %+v", task)), nil
}

func getProjectsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := getVikunjaClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	projects, err := client.GetProjects()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get projects: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Projects: %+v", projects)), nil
}

func createProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := getVikunjaClient()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	title, err := request.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError("title is required and must be a string"), nil
	}

	description := request.GetString("description", "")

	project, err := client.CreateProject(title, description)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create project: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Created project: %+v", project)), nil
}
