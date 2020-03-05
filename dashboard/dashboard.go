package dashboard

import (
	"fmt"
	"sort"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
)

type Monitor struct {
	Name     string
	Workflow string
	Branch   string
	Status   string
	Link     string
}

type Monitors []Monitor

func NewMonitor(project circleci.Project, pipeline circleci.Pipeline, workflow circleci.Workflow, status, link string) Monitor {
	return Monitor{
		Name:     project.Name(),
		Workflow: workflow.Name,
		Branch:   pipeline.VCS.Branch,
		Status:   status,
		Link:     link,
	}
}

func Build(circleCIClient circleci.CircleCI, filter *circleci.Filter) (Monitors, error) {
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
		filteredPipelines := pipelines.FilteredPerBranch()
		for branch, pipeline := range pipelines.LatestPerBranch() {
			workflows, err := circleCIClient.GetWorkflowsForPipeline(pipeline)
			if err != nil {
				return nil, err
			}
			dashboardData, err = dashboardData.AddWorkflows(circleCIClient, project, pipeline, workflows, filteredPipelines[branch])
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

func (d Monitors) AddWorkflows(circleCIClient circleci.CircleCI, project circleci.Project, pipeline circleci.Pipeline, workflows circleci.Workflows, filteredPipelines circleci.Pipelines) (Monitors, error) {
	for _, workflow := range workflows {
		monitor := NewMonitor(project, pipeline, workflow, "", "")
		if d.AlreadyExists(monitor) {
			continue
		}
		var status string
		if workflow.Name == "Build Error" {
			status = "errored"
		} else {
			var err error
			status, err = circleCIClient.WorkflowStatus(filteredPipelines, workflow)
			if err != nil {
				return nil, err
			}
		}
		link := circleCIClient.WorkflowLink(project, pipeline, workflow)
		monitor.Status = status
		monitor.Link = link
		d = append(d, monitor)
	}
	return d, nil
}
