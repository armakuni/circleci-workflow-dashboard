package dashboard_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
	"github.com/armakuni/circleci-workflow-dashboard/dashboard"
	"github.com/armakuni/circleci-workflow-dashboard/mocks"
)

var _ = Describe("#NewMonitor", func() {
	var (
		project = circleci.Project{
			VCSType:  "github",
			Username: "foobar",
			Reponame: "example",
		}
		pipeline = circleci.Pipeline{
			ID:  "1",
			VCS: circleci.VCS{Branch: "master"},
		}
		workflow = circleci.Workflow{
			ID:   "1",
			Name: "test-workflow",
		}
		status = "success"
		link   = "https://foobar.com"
	)

	It("returns a properly formatted dashboard monitor", func() {
		Ω(dashboard.NewMonitor(project, pipeline, workflow, status, link)).Should(Equal(dashboard.Monitor{
			Name:     "foobar/example",
			Workflow: "test-workflow",
			Branch:   "master",
			Status:   status,
			Link:     link,
		}))
	})
})

var _ = Describe("Monitors", func() {
	Describe("#Sort", func() {
		It("sorts the monitors by name, workflow and branch", func() {
			var monitors = dashboard.Monitors{
				{
					Name:     "foobar/example",
					Workflow: "test-workflow",
					Branch:   "zzz",
				},
				{
					Name:     "foobar/example",
					Workflow: "test-workflow",
					Branch:   "master",
				},
				{
					Name:     "foobar/example",
					Workflow: "zzz",
					Branch:   "master",
				},
				{
					Name:     "another/example",
					Workflow: "zzz",
					Branch:   "master",
				},
			}
			monitors.Sort()
			Ω(monitors).Should(Equal(dashboard.Monitors{
				{
					Name:     "another/example",
					Workflow: "zzz",
					Branch:   "master",
				},
				{
					Name:     "foobar/example",
					Workflow: "test-workflow",
					Branch:   "master",
				},
				{
					Name:     "foobar/example",
					Workflow: "test-workflow",
					Branch:   "zzz",
				},
				{
					Name:     "foobar/example",
					Workflow: "zzz",
					Branch:   "master",
				},
			}))
		})
	})

	Describe("#AlreadyExists", func() {
		var monitors = dashboard.Monitors{
			{
				Name:     "foobar/example",
				Workflow: "test-workflow",
				Branch:   "master",
			},
		}

		Context("when a dashboard monitor is already in monitors", func() {
			var monitor = dashboard.Monitor{
				Name:     "foobar/example",
				Workflow: "test-workflow",
				Branch:   "master",
			}

			It("returns true", func() {
				Ω(monitors.AlreadyExists(monitor)).Should(BeTrue())
			})
		})

		Context("when a dashboard monitor is new", func() {
			var monitor = dashboard.Monitor{
				Name:     "foobar/example",
				Workflow: "test-workflow",
				Branch:   "develop",
			}

			It("returns false", func() {
				Ω(monitors.AlreadyExists(monitor)).Should(BeFalse())
			})
		})
	})

	Describe("#AddWorkflows", func() {
		var (
			circleCIClient = &mocks.CircleCI{}
			monitors       dashboard.Monitors
			project        = circleci.Project{
				VCSType:  "github",
				Username: "foobar",
				Reponame: "example",
			}
			pipeline = circleci.Pipeline{
				ID:  "1",
				VCS: circleci.VCS{Branch: "master"},
			}
			workflows = circleci.Workflows{
				{
					ID:   "1",
					Name: "test-workflow",
				},
				{
					ID:   "1",
					Name: "test-workflow",
				},
			}
			filteredPipelines = circleci.Pipelines{pipeline}
		)

		AfterEach(func() {
			monitors = dashboard.Monitors{}
			circleCIClient = &mocks.CircleCI{}
		})

		Context("when getting workflow status errors", func() {
			BeforeEach(func() {
				circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("", fmt.Errorf("Error getting status"))
			})

			It("returns an error", func() {
				monitors, err := monitors.AddWorkflows(circleCIClient, project, pipeline, workflows, filteredPipelines)
				Ω(err).Should(MatchError("Error getting status"))
				Ω(monitors).Should(HaveLen(0))
			})
		})

		Context("when getting workflow status is successful", func() {
			BeforeEach(func() {
				circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("success", nil)
				circleCIClient.On("WorkflowLink", mock.Anything, mock.Anything, mock.Anything).Return("https://foobar.com")
			})

			It("adds a monitor per new workflow", func() {
				monitors, err := monitors.AddWorkflows(circleCIClient, project, pipeline, workflows, filteredPipelines)
				Ω(err).Should(BeNil())
				Ω(monitors).Should(Equal(dashboard.Monitors{{
					Name:     "foobar/example",
					Workflow: "test-workflow",
					Branch:   "master",
					Status:   "success",
					Link:     "https://foobar.com",
				}}))
			})
		})
	})
})

var _ = Describe("#Build", func() {
	var (
		circleCIClient = &mocks.CircleCI{}
		project        = circleci.Project{
			VCSType:  "github",
			Username: "foobar",
			Reponame: "example",
			Branches: map[string]interface{}{
				"master": nil,
			},
		}
		projects = circleci.Projects{project}
		pipeline = circleci.Pipeline{
			ID:  "1",
			VCS: circleci.VCS{Branch: "master"},
		}
		workflows = circleci.Workflows{
			{
				ID:   "1",
				Name: "test-workflow",
			},
			{
				ID:   "1",
				Name: "test-workflow",
			},
		}
		filteredPipelines = circleci.Pipelines{pipeline}
		filter            = circleci.Filter{}
	)

	AfterEach(func() {
		circleCIClient = &mocks.CircleCI{}
	})

	Context("when getting projects errors", func() {
		BeforeEach(func() {
			circleCIClient.On("GetAllProjects").Return(nil, fmt.Errorf("Error getting projects"))
		})

		It("returns an error", func() {
			monitors, err := dashboard.Build(circleCIClient, &filter)
			Ω(err).Should(MatchError("Error getting projects"))
			Ω(monitors).Should(BeNil())
		})
	})

	Context("when getting projects is successful", func() {
		BeforeEach(func() {
			circleCIClient.On("GetAllProjects").Return(projects, nil)
		})

		Context("and getting pipelines errors", func() {
			BeforeEach(func() {
				circleCIClient.On("GetAllPipelines", project).Return(nil, fmt.Errorf("Error getting pipelines"))
			})

			It("returns an error", func() {
				monitors, err := dashboard.Build(circleCIClient, &filter)
				Ω(err).Should(MatchError("Error getting pipelines"))
				Ω(monitors).Should(BeNil())
			})
		})

		Context("and gettings pipelines is successful", func() {
			BeforeEach(func() {
				circleCIClient.On("GetAllPipelines", project).Return(filteredPipelines, nil)
			})

			Context("and getting workflows errors", func() {
				BeforeEach(func() {
					circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(nil, fmt.Errorf("Error getting workflows"))
				})

				It("returns an error", func() {
					monitors, err := dashboard.Build(circleCIClient, &filter)
					Ω(err).Should(MatchError("Error getting workflows"))
					Ω(monitors).Should(BeNil())
				})
			})

			Context("and getting workflows is successful", func() {
				BeforeEach(func() {
					circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(workflows, nil)
				})

				Context("and adding workflows errors", func() {
					BeforeEach(func() {
						circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("", fmt.Errorf("Error adding workflows"))
					})

					It("returns an error", func() {
						monitors, err := dashboard.Build(circleCIClient, &filter)
						Ω(err).Should(MatchError("Error adding workflows"))
						Ω(monitors).Should(BeNil())
					})
				})

				Context("and adding workflows is successful", func() {
					BeforeEach(func() {
						circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("success", nil)
						circleCIClient.On("WorkflowLink", mock.Anything, mock.Anything, mock.Anything).Return("https://foobar.com")
					})

					It("returns dashboard monitors", func() {
						monitors, err := dashboard.Build(circleCIClient, &filter)
						Ω(err).Should(BeNil())
						Ω(monitors).Should(Equal(dashboard.Monitors{{
							Name:     "foobar/example",
							Workflow: "test-workflow",
							Branch:   "master",
							Status:   "success",
							Link:     "https://foobar.com",
						}}))
					})
				})
			})
		})
	})
})
