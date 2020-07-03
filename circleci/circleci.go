package circleci

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	statusSuccess      = "success"
	statusRunning      = "running"
	statusNotRunn      = "not_run"
	statusFailed       = "failed"
	statusError        = "error"
	statusFailing      = "failing"
	statusOnHold       = "on_hold"
	statusCanceled     = "canceled"
	statusUnauthorized = "unauthorized"
	statusUnknown      = "unknown"
)

const statusRespError = "Status Code was not ok"

var completedStauses = map[string]interface{}{
	statusSuccess: nil,
	statusFailed:  nil,
	statusError:   nil,
	statusUnknown: nil,
}

type PagedResponse struct {
	Items         json.RawMessage `json:"items"`
	NextPageToken *string         `json:"next_page_token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type CircleCI interface {
	GetAllProjects() (Projects, error)
	GetProjectEnvVars(string) (ProjectEnvVars, error)
	CreateProjectEnvVar(string, string, string) (ProjectEnvVar, error)
	DeleteProjectEnvVar(string, string) error
	GetAllPipelines(Project) (Pipelines, error)
	GetWorkflowsForPipeline(Pipeline) (Workflows, error)
	GetJobsForWorkflow(Workflow) (Jobs, error)
	PreviousCompleteWorkflowState(Pipelines, string) (string, error)
	WorkflowLink(Project, Pipeline, Workflow) string
	JobLink(Project, Job) string
	WorkflowStatus(Pipelines, Workflow) (string, error)
}

type Client struct {
	Config *Config
	Client *resty.Client
}

type Config struct {
	APIURL   string
	JobsURL  string
	APIToken string
}

func DefaultConfig() *Config {
	return &Config{
		APIURL:  "https://circleci.com",
		JobsURL: "https://app.circleci.com",
	}
}

func NewClient(config *Config) (*Client, error) {
	defaultConfig := DefaultConfig()
	if config.APIURL == "" {
		config.APIURL = defaultConfig.APIURL
	}
	if config.JobsURL == "" {
		config.JobsURL = defaultConfig.JobsURL
	}
	if config.APIToken == "" {
		return nil, fmt.Errorf("Must provide an API Token")
	}
	client := resty.New()
	client.SetBasicAuth(config.APIToken, "")
	return &Client{Client: client, Config: config}, nil
}

func (c *Client) get(urlPath string) (*resty.Response, error) {
	return c.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/%s", c.Config.APIURL, urlPath))
}

func (c *Client) post(urlPath string, body interface{}) (*resty.Response, error) {
	return c.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("%s/%s", c.Config.APIURL, urlPath))
}

func (c *Client) delete(urlPath string) (*resty.Response, error) {
	return c.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Delete(fmt.Sprintf("%s/%s", c.Config.APIURL, urlPath))
}

func (c *Client) callAPIV1(apiTarget string) (*resty.Response, error) {
	return c.get(fmt.Sprintf("api/v1.1/%s", apiTarget))
}

func (c *Client) callAPIV2(apiTarget string) (*resty.Response, error) {
	return c.get(fmt.Sprintf("api/v2/%s", apiTarget))
}

func (c *Client) pagedCallAPIV2(apiTarget string) ([]json.RawMessage, error) {
	var (
		nextPage   string
		items      []json.RawMessage
		pageTarget string
	)

	nextPageToken := &nextPage
	for nextPageToken != nil {
		if *nextPageToken == "" {
			pageTarget = apiTarget
		} else {
			query_joiner := "?"
			if strings.ContainsRune(apiTarget, '?') {
				query_joiner = "&"
			}
			pageTarget = fmt.Sprintf("%s%spage-token=%s", apiTarget, query_joiner, *nextPageToken)
		}
		resp, err := c.callAPIV2(pageTarget)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() > 299 {
			return nil, fmt.Errorf(statusRespError)
		}
		var pagedResponse PagedResponse
		if err := json.Unmarshal(resp.Body(), &pagedResponse); err != nil {
			return nil, err
		}
		items = append(items, pagedResponse.Items)
		nextPageToken = pagedResponse.NextPageToken
	}
	return items, nil
}

func (c *Client) GetAllProjects() (Projects, error) {
	resp, err := c.callAPIV1("projects")
	if err != nil {
		return nil, err
	}
	var projects Projects
	err = json.Unmarshal(resp.Body(), &projects)
	return projects, err
}

func (c *Client) GetProjectEnvVars(projectSlug string) (ProjectEnvVars, error) {
	var envVars ProjectEnvVars
	items, err := c.pagedCallAPIV2(fmt.Sprintf("project/%s/envvar", projectSlug))
	if err != nil {
		if err.Error() == statusRespError {
			return nil, fmt.Errorf("Could not get project, do you have the correct projectSlug and API permissions?")
		}
		return nil, err
	}
	for _, item := range items {
		if len(item) == 0 {
			continue
		}
		var pagedEnvVars ProjectEnvVars
		if err := json.Unmarshal(item, &pagedEnvVars); err != nil {
			return nil, err
		}
		envVars = append(envVars, pagedEnvVars...)
	}
	return envVars, nil
}

func (c *Client) CreateProjectEnvVar(projectSlug, key, value string) (ProjectEnvVar, error) {
	envVar := ProjectEnvVar{
		Name:  key,
		Value: value,
	}
	resp, err := c.post(fmt.Sprintf("api/v2/project/%s/envvar", projectSlug), envVar)
	if err != nil {
		return ProjectEnvVar{}, err
	}
	var respEnvVar ProjectEnvVar
	err = json.Unmarshal(resp.Body(), &respEnvVar)
	return respEnvVar, err
}

func (c *Client) DeleteProjectEnvVar(projectSlug, key string) error {
	_, err := c.delete(fmt.Sprintf("api/v2/project/%s/envvar/%s", projectSlug, key))
	return err
}

func (c *Client) GetAllPipelines(project Project) (Pipelines, error) {
	var pipelines Pipelines
	for branch, _ := range project.Branches {
		items, err := c.pagedCallAPIV2(fmt.Sprintf("project/%s/pipeline?branch=%s", project.Slug(), branch))
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			if len(item) == 0 {
				continue
			}
			var pagedPipelines Pipelines
			if err := json.Unmarshal(item, &pagedPipelines); err != nil {
				return nil, err
			}
			pipelines = append(pipelines, pagedPipelines...)
		}
	}
	return pipelines, nil
}

func (c *Client) GetWorkflowsForPipeline(pipeline Pipeline) (Workflows, error) {
	var workflows Workflows
	items, err := c.pagedCallAPIV2(fmt.Sprintf("pipeline/%s/workflow", pipeline.ID))
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if len(item) == 0 {
			continue
		}
		var pagedWorkflows Workflows
		if err := json.Unmarshal(item, &pagedWorkflows); err != nil {
			return nil, err
		}
		workflows = append(workflows, pagedWorkflows...)
	}
	return workflows, nil
}

func (c *Client) GetJobsForWorkflow(workflow Workflow) (Jobs, error) {
	var jobs Jobs
	items, err := c.pagedCallAPIV2(fmt.Sprintf("workflow/%s/job", workflow.ID))
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if len(item) == 0 {
			continue
		}
		var pagedJobs Jobs
		if err := json.Unmarshal(item, &pagedJobs); err != nil {
			return nil, err
		}
		jobs = append(jobs, pagedJobs...)
	}
	return jobs, nil
}

func (c *Client) PreviousCompleteWorkflowState(pipelines Pipelines, workflowName string) (string, error) {
	status := statusUnknown
	for _, pipeline := range pipelines {
		var workflowInPipeline bool
		workflows, err := c.GetWorkflowsForPipeline(pipeline)
		if err != nil {
			return "", err
		}
		for _, workflow := range workflows {
			if workflow.Name == "Build Error" {
				workflowInPipeline = true
			}
			if workflow.Name == workflowName {
				workflowInPipeline = true
				if _, ok := completedStauses[workflow.Status]; ok {
					return workflow.Status, nil
				}
			}
		}
		if !workflowInPipeline {
			return status, nil
		}
	}
	return status, nil
}

func (c *Client) WorkflowLink(project Project, pipeline Pipeline, workflow Workflow) string {
	return fmt.Sprintf("%s/pipelines/%s/%d/workflows/%s", c.Config.JobsURL, project.Slug(), pipeline.Number, workflow.ID)
}

func (c *Client) JobLink(project Project, job Job) string {
	return fmt.Sprintf("%s/jobs/%s/%s", c.Config.JobsURL, project.Slug(), job.ID)
}

func (c *Client) WorkflowStatus(pipelines Pipelines, workflow Workflow) (string, error) {
	status := workflow.Status
	if _, ok := completedStauses[workflow.Status]; !ok {
		previous_status, err := c.PreviousCompleteWorkflowState(pipelines, workflow.Name)
		if err != nil {
			return "", err
		}
		status = fmt.Sprintf("%s %s", status, previous_status)
	}
	return status, nil
}
