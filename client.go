package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type VikunjaClient struct {
	baseURL string
	token   string
	client  *http.Client
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	ProjectID   int    `json:"project_id"`
}

type Project struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewVikunjaClient(baseURL, token string) *VikunjaClient {
	return &VikunjaClient{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}
}

func getVikunjaClient() (*VikunjaClient, error) {
	godotenv.Load()

	baseURL := os.Getenv("VIKUNJA_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("VIKUNJA_URL environment variable is required")
	}

	token := os.Getenv("VIKUNJA_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VIKUNJA_TOKEN environment variable is required")
	}

	return NewVikunjaClient(baseURL, token), nil
}

func (c *VikunjaClient) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + "/api/v1" + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}

func (c *VikunjaClient) GetTasks() ([]Task, error) {
	path := "/tasks/all?filter=done%20%3D%20false%20%26%26%20due_date%20%3C%20now%20%26%26%20project%20%21%3D%2026&" +
		"sort_by=due_date&order_by=asc&sort_by=id&order_by=desc&page=1"
	resp, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *VikunjaClient) CreateTask(projectID int, title string) (*Task, error) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	taskData := map[string]interface{}{
		"title":    title,
		"due_date": now,
		"asignees": []int{1},
	}

	resp, err := c.makeRequest("PUT", fmt.Sprintf("/projects/%d/tasks", projectID), taskData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *VikunjaClient) UpdateTask(taskID int, updates map[string]interface{}) (*Task, error) {
	resp, err := c.makeRequest("POST", fmt.Sprintf("/tasks/%d", taskID), updates)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *VikunjaClient) GetProjects() ([]Project, error) {
	resp, err := c.makeRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *VikunjaClient) CreateProject(title, description string) (*Project, error) {
	projectData := map[string]interface{}{
		"title":       title,
		"description": description,
	}

	resp, err := c.makeRequest("PUT", "/projects", projectData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}
