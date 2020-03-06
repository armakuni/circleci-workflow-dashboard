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
			filteredPipelines  = circleci.Pipelines{pipeline}
			buildError         = false
			animatedBuildError = true
			workflowInfo       dashboard.WorkflowDetails
		)

		BeforeEach(func() {
			workflowInfo = dashboard.WorkflowDetails{
				Workflows:         workflows,
				Project:           project,
				Pipeline:          pipeline,
				FilteredPipelines: filteredPipelines,
				BuildError:        buildError,
			}
		})

		JustBeforeEach(func() {
			workflowInfo.CircleCIClient = circleCIClient
		})

		AfterEach(func() {
			monitors = dashboard.Monitors{}
			circleCIClient = &mocks.CircleCI{}
		})

		Context("when getting workflow status errors", func() {
			BeforeEach(func() {
				circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("", fmt.Errorf("Error getting status"))
			})

			It("returns an error", func() {
				monitors, err := monitors.AddWorkflows(workflowInfo, animatedBuildError)
				Ω(err).Should(MatchError("Error getting status"))
				Ω(monitors).Should(HaveLen(0))
			})
		})

		Context("when getting workflow status is successful", func() {
			BeforeEach(func() {
				circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("success", nil)
				circleCIClient.On("WorkflowLink", mock.Anything, mock.Anything, mock.Anything).Return("https://foobar.com")
			})

			Context("and there was a build error", func() {
				BeforeEach(func() {
					workflowInfo.BuildError = true
				})

				Context("and animated build is true", func() {
					It("adds a monitor per new workflow", func() {
						monitors, err := monitors.AddWorkflows(workflowInfo, animatedBuildError)
						Ω(err).Should(BeNil())
						Ω(monitors).Should(Equal(dashboard.Monitors{{
							Name:     "foobar/example",
							Workflow: "test-workflow",
							Branch:   "master",
							Status:   "success errored",
							Link:     "https://foobar.com",
						}}))
					})
				})

				Context("and animated build is false", func() {
					It("adds a monitor per new workflow", func() {
						monitors, err := monitors.AddWorkflows(workflowInfo, false)
						Ω(err).Should(BeNil())
						Ω(monitors).Should(Equal(dashboard.Monitors{{
							Name:     "foobar/example",
							Workflow: "test-workflow",
							Branch:   "master",
							Status:   "success errored-static",
							Link:     "https://foobar.com",
						}}))
					})
				})
			})

			Context("and there was not a build error", func() {
				It("adds a monitor per new workflow", func() {
					monitors, err := monitors.AddWorkflows(workflowInfo, animatedBuildError)
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
		pipeline2 = circleci.Pipeline{
			ID:  "2",
			VCS: circleci.VCS{Branch: "master"},
		}
		workflows = circleci.Workflows{
			{
				ID:   "1",
				Name: "Build Error",
			},
		}
		workflows2 = circleci.Workflows{
			{
				ID:   "2",
				Name: "test-workflow",
			},
		}
		filteredPipelines  = circleci.Pipelines{pipeline, pipeline2}
		filter             = circleci.Filter{}
		animatedBuildError = true
	)

	AfterEach(func() {
		circleCIClient = &mocks.CircleCI{}
	})

	Context("when getting projects errors", func() {
		BeforeEach(func() {
			circleCIClient.On("GetAllProjects").Return(nil, fmt.Errorf("Error getting projects"))
		})

		It("returns an error", func() {
			monitors, err := dashboard.Build(circleCIClient, &filter, animatedBuildError)
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
				monitors, err := dashboard.Build(circleCIClient, &filter, animatedBuildError)
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
					monitors, err := dashboard.Build(circleCIClient, &filter, animatedBuildError)
					Ω(err).Should(MatchError("Error getting workflows"))
					Ω(monitors).Should(BeNil())
				})
			})

			Context("and getting workflows is successful", func() {
				Context("and getting previous workflows without build errors is successful", func() {
					BeforeEach(func() {
						circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(workflows, nil)
						circleCIClient.On("GetWorkflowsForPipeline", pipeline2).Return(nil, fmt.Errorf("Error getting previous workflows"))
					})

					It("returns an error", func() {
						monitors, err := dashboard.Build(circleCIClient, &filter, animatedBuildError)
						Ω(err).Should(MatchError("Error getting previous workflows"))
						Ω(monitors).Should(BeNil())
					})
				})

				Context("and getting previous workflows without build error is successful", func() {
					BeforeEach(func() {
						circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(workflows, nil)
						circleCIClient.On("GetWorkflowsForPipeline", pipeline2).Return(workflows2, nil)
					})

					Context("and adding workflows errors", func() {
						BeforeEach(func() {
							circleCIClient.On("WorkflowStatus", mock.Anything, mock.Anything).Return("", fmt.Errorf("Error adding workflows"))
						})

						It("returns an error", func() {
							monitors, err := dashboard.Build(circleCIClient, &filter, animatedBuildError)
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
							monitors, err := dashboard.Build(circleCIClient, &filter, animatedBuildError)
							Ω(err).Should(BeNil())
							Ω(monitors).Should(Equal(dashboard.Monitors{{
								Name:     "foobar/example",
								Workflow: "test-workflow",
								Branch:   "master",
								Status:   "success errored",
								Link:     "https://foobar.com",
							}}))
						})
					})
				})
			})
		})
	})
})

var _ = Describe("WorkflowDetails", func() {
	Describe("#GetLatestWorkflowWithoutBuildError", func() {
		var (
			circleCIClient = &mocks.CircleCI{}
			pipeline       = circleci.Pipeline{
				ID:  "1",
				VCS: circleci.VCS{Branch: "master"},
			}
			pipeline2 = circleci.Pipeline{
				ID:  "2",
				VCS: circleci.VCS{Branch: "master"},
			}
			workflows = circleci.Workflows{
				{
					ID:   "1",
					Name: "Build Error",
				},
			}
			workflows2 = circleci.Workflows{
				{
					ID:   "2",
					Name: "test-workflow",
				},
			}
			filteredPipelines = circleci.Pipelines{pipeline, pipeline2}
			workflowDetails   dashboard.WorkflowDetails
		)

		BeforeEach(func() {
			workflowDetails = dashboard.WorkflowDetails{}
		})

		JustBeforeEach(func() {
			workflowDetails.CircleCIClient = circleCIClient
		})

		AfterEach(func() {
			circleCIClient = &mocks.CircleCI{}
		})

		Context("when the workflows do not contain a build error", func() {
			JustBeforeEach(func() {
				workflowDetails.Workflows = workflows2
				workflowDetails.FilteredPipelines = filteredPipelines
			})

			It("returns the workflows and pipelines as they are", func() {
				err := workflowDetails.GetLatestWorkflowWithoutBuildError()
				Ω(workflowDetails.Workflows).Should(Equal(workflows2))
				Ω(workflowDetails.Pipeline).Should(Equal(pipeline))
				Ω(workflowDetails.FilteredPipelines).Should(Equal(filteredPipelines))
				Ω(workflowDetails.BuildError).Should(BeFalse())
				Ω(err).Should(BeNil())
			})
		})

		Context("when the workflows contain a build error", func() {
			Context("when getting workflows errors", func() {
				JustBeforeEach(func() {
					workflowDetails.Workflows = workflows
					workflowDetails.FilteredPipelines = filteredPipelines
				})

				BeforeEach(func() {
					circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(nil, fmt.Errorf("Getting workflows errored"))
					circleCIClient.On("GetWorkflowsForPipeline", pipeline2).Return(workflows2, nil)
				})

				It("returns an error", func() {
					err := workflowDetails.GetLatestWorkflowWithoutBuildError()
					Ω(err).Should(MatchError("Getting workflows errored"))
				})
			})

			Context("when getting workflows is successful", func() {
				JustBeforeEach(func() {
					workflowDetails.Workflows = workflows
					workflowDetails.FilteredPipelines = filteredPipelines
				})

				Context("when there is a previous working workflow", func() {
					BeforeEach(func() {
						circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(workflows, nil)
						circleCIClient.On("GetWorkflowsForPipeline", pipeline2).Return(workflows2, nil)
					})

					It("returns the last good workflow and pipelines", func() {
						err := workflowDetails.GetLatestWorkflowWithoutBuildError()
						Ω(workflowDetails.Workflows).Should(Equal(workflows2))
						Ω(workflowDetails.Pipeline).Should(Equal(pipeline2))
						Ω(workflowDetails.FilteredPipelines).Should(Equal(circleci.Pipelines{pipeline2}))
						Ω(workflowDetails.BuildError).Should(BeTrue())
						Ω(err).Should(BeNil())
					})
				})

				Context("when the workflows have only ever been build errors", func() {
					BeforeEach(func() {
						circleCIClient.On("GetWorkflowsForPipeline", pipeline).Return(workflows, nil)
						circleCIClient.On("GetWorkflowsForPipeline", pipeline2).Return(workflows, nil)
					})

					It("returns the earliest existance of the build error with an unknown status", func() {
						err := workflowDetails.GetLatestWorkflowWithoutBuildError()
						Ω(workflowDetails.Workflows).Should(Equal(circleci.Workflows{{ID: "1", Name: "Build Error", Status: "unknown"}}))
						Ω(workflowDetails.Pipeline).Should(Equal(pipeline2))
						Ω(workflowDetails.FilteredPipelines).Should(Equal(circleci.Pipelines{pipeline2}))
						Ω(workflowDetails.BuildError).Should(BeTrue())
						Ω(err).Should(BeNil())
					})
				})
			})
		})
	})
})
