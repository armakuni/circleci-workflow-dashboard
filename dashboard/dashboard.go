package dashboard

import (
	"fmt"
	"sort"
	"strings"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
)

type FeatureFlags struct {
	AnimatedBuildErrors bool
}

type Monitor struct {
	Name     string
	Workflow string
	Branch   string
	Status   string
	Link     string
}

type MonitorConfig struct {
	HideOrganization bool
	HideBranch       bool
	BranchFilter     string
}

type Monitors []Monitor

func NewMonitor(project circleci.Project, pipeline circleci.Pipeline, workflow circleci.Workflow, status, link string, config *MonitorConfig) Monitor {
	projectName := project.Name()
	branchName := pipeline.VCS.Branch

	if config.HideOrganization {
		projectName = strings.Join(strings.Split(projectName, "/")[1:], "/")
	}

	if config.HideBranch {
		branchName = ""
	}

	return Monitor{
		Name:     projectName,
		Workflow: workflow.Name,
		Branch:   branchName,
		Status:   status,
		Link:     link,
	}
}

func Build(circleCIClient circleci.CircleCI, filter *circleci.Filter, featureFlags *FeatureFlags, monitorConfig *MonitorConfig) (Monitors, error) {
	var dashboardData Monitors
	projects, err := circleCIClient.GetAllProjects()
	if err != nil {
		return nil, err
	}
	projects = projects.Filter(filter)
	for _, project := range projects {
		pipelines, err := circleCIClient.GetAllPipelines(project)
		if err != nil {
			return nil, err
		}
		filteredPipelines := pipelines.FilteredPerBranch(monitorConfig.BranchFilter)
		for branch, pipeline := range pipelines.LatestPerBranch() {
			workflows, err := circleCIClient.GetWorkflowsForPipeline(pipeline)
			if err != nil {
				return nil, err
			}
			workflowInfo := WorkflowDetails{
				Project:           project,
				CircleCIClient:    circleCIClient,
				Workflows:         workflows,
				FilteredPipelines: filteredPipelines[branch],
			}
			err = workflowInfo.GetLatestWorkflowWithoutBuildError()
			if err != nil {
				return nil, err
			}
			dashboardData, err = dashboardData.AddWorkflows(workflowInfo, featureFlags, monitorConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	dashboardData.Sort()
	return dashboardData, nil
}

func (d *Monitors) Sort() {
	monitors := *d
	sort.Slice(monitors, func(i, j int) bool {
		iMonitor := fmt.Sprintf("%s-%s-%s", monitors[i].Name, monitors[i].Workflow, monitors[i].Branch)
		jMonitor := fmt.Sprintf("%s-%s-%s", monitors[j].Name, monitors[j].Workflow, monitors[j].Branch)
		return iMonitor < jMonitor
	})
	d = &monitors
}

func (d *Monitors) AlreadyExists(monitor Monitor) bool {
	for _, mon := range *d {
		if mon.Name == monitor.Name &&
			mon.Workflow == monitor.Workflow &&
			mon.Branch == monitor.Branch {
			return true
		}
	}
	return false
}

func (d Monitors) AddWorkflows(workflowInfo WorkflowDetails, featureFlags *FeatureFlags, monitorConfig *MonitorConfig) (Monitors, error) {
	for _, workflow := range workflowInfo.Workflows {
		monitor := NewMonitor(workflowInfo.Project, workflowInfo.Pipeline, workflow, "", "", monitorConfig)
		if d.AlreadyExists(monitor) {
			continue
		}
		status, err := workflowInfo.CircleCIClient.WorkflowStatus(workflowInfo.FilteredPipelines, workflow)
		if err != nil {
			return nil, err
		}
		if workflowInfo.BuildError {
			errorStatus := "errored"
			if !featureFlags.AnimatedBuildErrors {
				errorStatus = "errored-static"
			}
			status = fmt.Sprintf("%s %s", status, errorStatus)
		}
		link := workflowInfo.CircleCIClient.WorkflowLink(workflowInfo.Project, workflowInfo.Pipeline, workflow)
		monitor.Status = status
		monitor.Link = link
		d = append(d, monitor)
	}
	return d, nil
}

type WorkflowDetails struct {
	CircleCIClient    circleci.CircleCI
	Project           circleci.Project
	Workflows         circleci.Workflows
	Pipeline          circleci.Pipeline
	FilteredPipelines circleci.Pipelines
	BuildError        bool
}

func (workflowInfo *WorkflowDetails) GetLatestWorkflowWithoutBuildError() error {
	if workflowInfo.Workflows.BuildError() {
		workflowInfo.BuildError = true
		previousNonError, err := workflowInfo.getPreviousWorkflowDetails()
		if err != nil {
			return err
		}
		if previousNonError {
			return nil
		} else {
			workflow := workflowInfo.Workflows[0]
			workflow.Status = "unknown"
			workflowInfo.Workflows = circleci.Workflows{workflow}
			workflowInfo.Pipeline = workflowInfo.FilteredPipelines[len(workflowInfo.FilteredPipelines)-1]
			workflowInfo.FilteredPipelines = circleci.Pipelines{workflowInfo.Pipeline}
			return nil
		}
	}
	workflowInfo.Pipeline = workflowInfo.FilteredPipelines[0]
	return nil
}

func (workflowInfo *WorkflowDetails) getPreviousWorkflowDetails() (bool, error) {
	for index, pipeline := range workflowInfo.FilteredPipelines {
		workflows, err := workflowInfo.CircleCIClient.GetWorkflowsForPipeline(pipeline)
		if err != nil {
			return false, err
		}
		if !workflows.BuildError() {
			workflowInfo.Workflows = workflows
			workflowInfo.Pipeline = pipeline
			workflowInfo.FilteredPipelines = workflowInfo.FilteredPipelines[index:]
			return true, nil
		}
	}
	return false, nil
}
