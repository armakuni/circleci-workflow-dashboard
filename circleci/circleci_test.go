package circleci_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
)

var _ = Describe("Client", func() {
	var (
		config *circleci.Config
		client *circleci.Client
	)

	JustBeforeEach(func() {
		var err error
		host := ""
		if server != nil {
			host = Host(server)
		}
		config = &circleci.Config{
			APIURL:   host,
			APIToken: "fakeToken",
		}
		client, err = circleci.NewClient(config)
		client.Client.SetDisableWarn(true)
		Ω(err).Should(BeNil())
	})

	AfterEach(func() {
		teardown()
	})

	Describe("#NewClient", func() {
		Context("when you don't provide an API token", func() {
			It("returns an error", func() {
				circleConfig := &circleci.Config{}
				_, err := circleci.NewClient(circleConfig)
				Ω(err).Should(MatchError("Must provide an API Token"))
			})
		})

		Context("when you provide an API token", func() {
			Context("and you only provide an API token", func() {
				It("returns a client with the default URLs", func() {
					circleConfig := &circleci.Config{APIToken: "foobar"}
					circleClient, err := circleci.NewClient(circleConfig)
					Ω(err).Should(BeNil())
					Ω(circleClient.Config.APIToken).Should(Equal("foobar"))
					Ω(circleClient.Config.APIURL).Should(Equal("https://circleci.com"))
					Ω(circleClient.Config.JobsURL).Should(Equal("https://app.circleci.com"))
				})
			})

			Context("and you provide all the config", func() {
				It("returns a client without using the defaults", func() {
					circleConfig := &circleci.Config{APIToken: "foobar", APIURL: "https://foobar.com", JobsURL: "https://app.foobar.com"}
					circleClient, err := circleci.NewClient(circleConfig)
					Ω(err).Should(BeNil())
					Ω(circleClient.Config.APIToken).Should(Equal("foobar"))
					Ω(circleClient.Config.APIURL).Should(Equal("https://foobar.com"))
					Ω(circleClient.Config.JobsURL).Should(Equal("https://app.foobar.com"))
				})
			})
		})
	})

	Describe("#GetAllProjects", func() {
		Context("when circleci returns an error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v1.1/projects", "[]", 500, "", nil},
				}
				setupMultiple(mocks)
			})

			JustBeforeEach(func() {
				client.Config.APIURL = "broken"
			})

			It("returns an error", func() {
				projects, err := client.GetAllProjects()
				Ω(err).Should(MatchError(`Get "//%2Fbroken%2Fapi%2Fv1.1%2Fprojects/broken/api/v1.1/projects": unsupported protocol scheme ""`))
				Ω(projects).Should(BeEmpty())
			})
		})

		Context("when circleci returns invalid JSON", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v1.1/projects", "[}", 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				projects, err := client.GetAllProjects()
				Ω(err).Should(MatchError("invalid character '}' looking for beginning of value"))
				Ω(projects).Should(BeEmpty())
			})
		})

		Context("when circleci returns valid JSON", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v1.1/projects", project_resp, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("gets all projects from CircleCI", func() {
				projects, err := client.GetAllProjects()
				Ω(err).Should(BeNil())
				Ω(projects).Should(HaveLen(2))
				Ω(projects[0].VCSType).Should(Equal("github"))
				Ω(projects[0].Username).Should(Equal("foobar"))
				Ω(projects[0].Reponame).Should(Equal("example"))
			})
		})
	})

	Describe("#GetProjectEnvVars", func() {
		var projectSlug = "github/foobar/example"

		Context("When the project doesn't exist", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/project/github/foobar/example/envvar", `{"message": "Project not found"}`, 404, "", nil},
				}
				setupMultiple(mocks)
			})

			It("raises an error", func() {
				envVars, err := client.GetProjectEnvVars(projectSlug)
				Ω(err).Should(MatchError("Could not get project, do you have the correct projectSlug and API permissions?"))
				Ω(envVars).Should(BeNil())
			})
		})

		Context("when the project exists", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/project/github/foobar/example/envvar", projecteEnvVarResp, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns the env vars of a project", func() {
				envVars, err := client.GetProjectEnvVars(projectSlug)
				Ω(err).Should(BeNil())
				Ω(envVars).Should(Equal(circleci.ProjectEnvVars{
					{
						Name:  "ARTIFACTORY_PASSWORD",
						Value: "xxxxxxxx",
					},
					{
						Name:  "ARTIFACTORY_USER",
						Value: "xxxxxxxx",
					},
					{
						Name:  "CACHE_VERSION",
						Value: "xxxxxxxx",
					},
					{Name: "PROJECT_NAME", Value: "xxxxxxxx"},
					{Name: "STUNNEL_HOST", Value: "xxxxxxxx"},
					{Name: "STUNNEL_PSK", Value: "xxxxxxxx"},
				}))
			})
		})
	})

	Describe("#CreateProjectEnvVar", func() {
		var projectSlug = "github/foobar/example"

		BeforeEach(func() {
			mocks := []MockRoute{
				{"POST", "/api/v2/project/github/foobar/example/envvar", createEnvVarResp, 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("creates the env var in a project", func() {
			envVar, err := client.CreateProjectEnvVar(projectSlug, "foo", "bar")
			Ω(err).Should(BeNil())
			Ω(envVar).Should(Equal(circleci.ProjectEnvVar{Name: "foo", Value: "xxxxxxx"}))
		})
	})

	Describe("#DeleteProjectEnvVar", func() {
		var projectSlug = "github/foobar/example"

		BeforeEach(func() {
			mocks := []MockRoute{
				{"DELETE", "/api/v2/project/github/foobar/example/envvar/foo", "", 200, "", nil},
			}
			setupMultiple(mocks)
		})

		It("deletes the env var in a project", func() {
			err := client.DeleteProjectEnvVar(projectSlug, "foo")
			Ω(err).Should(BeNil())
		})
	})

	Describe("#GetAllPipelines", func() {
		var project = circleci.Project{
			VCSType:  "github",
			Username: "foobar",
			Reponame: "example",
			Branches: map[string]interface{}{
				"master": nil,
			},
		}

		Context("when circleci returns an error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/project/github/foobar/example/pipeline", "[]", 500, "branch=master", nil},
				}
				setupMultiple(mocks)
			})

			JustBeforeEach(func() {
				client.Config.APIURL = "broken"
			})

			It("returns an error", func() {
				pipelines, err := client.GetAllPipelines(project)
				Ω(err).Should(MatchError(`Get "//%2Fbroken%2Fapi%2Fv2%2Fproject%2Fgithub%2Ffoobar%2Fexample%2Fpipeline%3Fbranch=master/broken/api/v2/project/github/foobar/example/pipeline?branch=master": unsupported protocol scheme ""`))
				Ω(pipelines).Should(BeEmpty())
			})
		})

		Context("when circlci returns invalid json", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/project/github/foobar/example/pipeline", "[}", 200, "branch=master", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				pipelines, err := client.GetAllPipelines(project)
				Ω(err).Should(MatchError("invalid character '}' looking for beginning of value"))
				Ω(pipelines).Should(BeEmpty())
			})
		})

		Context("when circleci returns valid json but invalid item types", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/project/github/foobar/example/pipeline", invalid_items_resp, 200, "branch=master", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				pipelines, err := client.GetAllPipelines(project)
				Ω(err).Should(MatchError("json: cannot unmarshal object into Go value of type circleci.Pipelines"))
				Ω(pipelines).Should(BeEmpty())
			})
		})

		Context("when circle return valid json", func() {
			Context("and there is a single page of pipelines", func() {
				BeforeEach(func() {
					mocks := []MockRoute{
						{"GET", "/api/v2/project/github/foobar/example/pipeline", pipeline_resp_single_page, 200, "branch=master", nil},
					}
					setupMultiple(mocks)
				})

				It("returns all pipelines", func() {
					pipelines, err := client.GetAllPipelines(project)
					Ω(err).Should(BeNil())
					Ω(pipelines).Should(HaveLen(2))
					Ω(pipelines[0].ID).Should(Equal("1"))
					Ω(pipelines[0].VCS.Branch).Should(Equal("master"))
				})
			})

			Context("and there are multiple pages of pipelines", func() {
				BeforeEach(func() {
					mocks := []MockRoute{
						{"GET", "/api/v2/project/github/foobar/example/pipeline", pipeline_resp_3, 200, "branch=master&page-token=page3", nil},
						{"GET", "/api/v2/project/github/foobar/example/pipeline", pipeline_resp_2, 200, "branch=master&page-token=page2", nil},
						{"GET", "/api/v2/project/github/foobar/example/pipeline", pipeline_resp_1, 200, "branch=master", nil},
					}
					setupMultiple(mocks)
				})

				It("returns all pipelines", func() {
					pipelines, err := client.GetAllPipelines(project)
					Ω(err).Should(BeNil())
					Ω(pipelines).Should(HaveLen(4))
					Ω(pipelines[3].ID).Should(Equal("4"))
					Ω(pipelines[3].VCS.Branch).Should(Equal("develop"))
				})
			})
		})
	})

	Describe("#GetWorkflowsForPipeline", func() {
		var pipeline = circleci.Pipeline{
			ID: "1",
			VCS: circleci.VCS{
				Branch: "master",
			},
		}

		Context("when circleci returns an error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", "[]", 500, "", nil},
				}
				setupMultiple(mocks)
			})

			JustBeforeEach(func() {
				client.Config.APIURL = "broken"
			})

			It("returns an error", func() {
				workflows, err := client.GetWorkflowsForPipeline(pipeline)
				Ω(err).Should(MatchError(`Get "//%2Fbroken%2Fapi%2Fv2%2Fpipeline%2F1%2Fworkflow/broken/api/v2/pipeline/1/workflow": unsupported protocol scheme ""`))
				Ω(workflows).Should(BeEmpty())
			})
		})

		Context("when circlci returns invalid json", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", "[}", 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				workflows, err := client.GetWorkflowsForPipeline(pipeline)
				Ω(err).Should(MatchError("invalid character '}' looking for beginning of value"))
				Ω(workflows).Should(BeEmpty())
			})
		})

		Context("when circleci returns valid json but invalid item types", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", invalid_items_resp, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				workflows, err := client.GetWorkflowsForPipeline(pipeline)
				Ω(err).Should(MatchError("json: cannot unmarshal object into Go value of type circleci.Workflows"))
				Ω(workflows).Should(BeEmpty())
			})
		})

		Context("when circle return valid json", func() {
			Context("and there is a single page of workflows", func() {
				BeforeEach(func() {
					mocks := []MockRoute{
						{"GET", "/api/v2/pipeline/1/workflow", workflows_resp_single_page, 200, "", nil},
					}
					setupMultiple(mocks)
				})

				It("returns all workflows", func() {
					workflows, err := client.GetWorkflowsForPipeline(pipeline)
					Ω(err).Should(BeNil())
					Ω(workflows).Should(HaveLen(2))
					Ω(workflows[0].ID).Should(Equal("1"))
					Ω(workflows[0].Name).Should(Equal("workflow1"))
					Ω(workflows[0].Status).Should(Equal("success"))
				})
			})

			Context("and there are multiple pages of workflows", func() {
				BeforeEach(func() {
					mocks := []MockRoute{
						{"GET", "/api/v2/pipeline/1/workflow", workflows_resp_3, 200, "page-token=page3", nil},
						{"GET", "/api/v2/pipeline/1/workflow", workflows_resp_2, 200, "page-token=page2", nil},
						{"GET", "/api/v2/pipeline/1/workflow", workflows_resp_1, 200, "", nil},
					}
					setupMultiple(mocks)
				})

				It("returns all workflows", func() {
					workflows, err := client.GetWorkflowsForPipeline(pipeline)
					Ω(err).Should(BeNil())
					Ω(workflows).Should(HaveLen(4))
					Ω(workflows[3].ID).Should(Equal("4"))
					Ω(workflows[3].Name).Should(Equal("workflow4"))
					Ω(workflows[3].Status).Should(Equal("failed"))
				})
			})
		})
	})

	Describe("#GetJobsForWorkflow", func() {
		var workflow = circleci.Workflow{
			ID: "1",
		}

		Context("when circleci returns an error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/workflow/1/job", "[]", 500, "", nil},
				}
				setupMultiple(mocks)
			})

			JustBeforeEach(func() {
				client.Config.APIURL = "broken"
			})

			It("returns an error", func() {
				jobs, err := client.GetJobsForWorkflow(workflow)
				Ω(err).Should(MatchError(`Get "//%2Fbroken%2Fapi%2Fv2%2Fworkflow%2F1%2Fjob/broken/api/v2/workflow/1/job": unsupported protocol scheme ""`))
				Ω(jobs).Should(BeEmpty())
			})
		})

		Context("when circlci returns invalid json", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/workflow/1/job", "[}", 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				jobs, err := client.GetJobsForWorkflow(workflow)
				Ω(err).Should(MatchError("invalid character '}' looking for beginning of value"))
				Ω(jobs).Should(BeEmpty())
			})
		})

		Context("when circleci returns valid json but invalid item types", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/workflow/1/job", invalid_items_resp, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				jobs, err := client.GetJobsForWorkflow(workflow)
				Ω(err).Should(MatchError("json: cannot unmarshal object into Go value of type circleci.Jobs"))
				Ω(jobs).Should(BeEmpty())
			})
		})

		Context("when circle return valid json", func() {
			Context("and there is a single page of jobs", func() {
				BeforeEach(func() {
					mocks := []MockRoute{
						{"GET", "/api/v2/workflow/1/job", jobs_resp_single_page, 200, "", nil},
					}
					setupMultiple(mocks)
				})

				It("returns all jobs", func() {
					jobs, err := client.GetJobsForWorkflow(workflow)
					Ω(err).Should(BeNil())
					Ω(jobs).Should(HaveLen(2))
					Ω(jobs[0].ID).Should(Equal("1"))
				})
			})

			Context("and there are multiple pages of jobs", func() {
				BeforeEach(func() {
					mocks := []MockRoute{
						{"GET", "/api/v2/workflow/1/job", jobs_resp_3, 200, "page-token=page3", nil},
						{"GET", "/api/v2/workflow/1/job", jobs_resp_2, 200, "page-token=page2", nil},
						{"GET", "/api/v2/workflow/1/job", jobs_resp_1, 200, "", nil},
					}
					setupMultiple(mocks)
				})

				It("returns all jobs", func() {
					jobs, err := client.GetJobsForWorkflow(workflow)
					Ω(err).Should(BeNil())
					Ω(jobs).Should(HaveLen(4))
					Ω(jobs[3].ID).Should(Equal("4"))
				})
			})
		})
	})

	Describe("#PreviousCompleteWorkflowState", func() {
		var pipelines = circleci.Pipelines{
			{
				ID: "1",
			},
			{
				ID: "2",
			},
			{
				ID: "3",
			},
		}

		Context("when getting workflows returns an error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", "[}", 200, "", nil},
					{"GET", "/api/v2/pipeline/2/workflow", workflows_resp_2, 200, "", nil},
					{"GET", "/api/v2/pipeline/3/workflow", workflows_resp_1, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				status, err := client.PreviousCompleteWorkflowState(pipelines, "previous_status")
				Ω(err).Should(MatchError("invalid character '}' looking for beginning of value"))
				Ω(status).To(Equal(""))
			})
		})

		Context("when status is known", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", workflow_resp_previous_state_1, 200, "", nil},
					{"GET", "/api/v2/pipeline/2/workflow", workflow_resp_previous_state_2, 200, "", nil},
					{"GET", "/api/v2/pipeline/3/workflow", workflow_resp_previous_state_3, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns the previous workflow status", func() {
				status, err := client.PreviousCompleteWorkflowState(pipelines, "previous_status")
				Ω(err).To(BeNil())
				Ω(status).To(Equal("success"))
			})
		})

		Context("when all builds are Build Error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", workflow_resp_previous_state_2, 200, "", nil},
					{"GET", "/api/v2/pipeline/2/workflow", workflow_resp_previous_state_2, 200, "", nil},
					{"GET", "/api/v2/pipeline/3/workflow", workflow_resp_previous_state_2, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns unknown for the status", func() {
				status, err := client.PreviousCompleteWorkflowState(pipelines, "previous_status")
				Ω(err).To(BeNil())
				Ω(status).To(Equal("unknown"))
			})
		})

		Context("when status is unknown", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", workflows_resp_3, 200, "", nil},
					{"GET", "/api/v2/pipeline/2/workflow", workflows_resp_3, 200, "", nil},
					{"GET", "/api/v2/pipeline/3/workflow", workflows_resp_3, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns unknown for the status", func() {
				status, err := client.PreviousCompleteWorkflowState(pipelines, "previous_status")
				Ω(err).To(BeNil())
				Ω(status).To(Equal("unknown"))
			})
		})
	})

	Describe("#WorkflowLink", func() {
		var (
			project = circleci.Project{
				VCSType:  "github",
				Username: "foobar",
				Reponame: "example",
			}
			pipeline = circleci.Pipeline{
				ID:     "1",
				Number: 36,
				VCS: circleci.VCS{
					Branch: "master",
				},
			}
			workflow = circleci.Workflow{
				ID: "2",
			}
		)

		It("returns a formatted workflow link to circleci", func() {
			workflowLink := client.WorkflowLink(project, pipeline, workflow)
			Ω(workflowLink).Should(Equal("https://app.circleci.com/pipelines/github/foobar/example/36/workflows/2"))
		})
	})

	Describe("#JobLink", func() {
		var (
			project = circleci.Project{
				VCSType:  "github",
				Username: "foobar",
				Reponame: "example",
			}
			job = circleci.Job{
				ID: "2",
			}
		)

		It("returns a formatted workflow link to circleci", func() {
			workflowLink := client.JobLink(project, job)
			Ω(workflowLink).Should(Equal("https://app.circleci.com/jobs/github/foobar/example/2"))
		})
	})

	Describe("#WorkflowStatus", func() {
		var (
			pipelines = circleci.Pipelines{
				{
					ID: "1",
				},
			}
			workflow = circleci.Workflow{
				ID:     "1",
				Status: "canceled",
				Name:   "previous_status",
			}
			completedWorkflow = circleci.Workflow{
				ID:     "1",
				Status: "failed",
				Name:   "foobar",
			}
		)
		Context("when circleci returns an error", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", "[}", 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("returns an error", func() {
				status, err := client.WorkflowStatus(pipelines, workflow)
				Ω(err).Should(MatchError("invalid character '}' looking for beginning of value"))
				Ω(status).To(Equal(""))
			})
		})

		Context("when the workflows status is complete", func() {
			It("returns the compelted status", func() {
				status, err := client.WorkflowStatus(pipelines, completedWorkflow)
				Ω(err).Should(BeNil())
				Ω(status).Should(Equal("failed"))
			})
		})

		Context("when the workflows status is not complete", func() {
			BeforeEach(func() {
				mocks := []MockRoute{
					{"GET", "/api/v2/pipeline/1/workflow", workflow_resp_previous_state_3, 200, "", nil},
				}
				setupMultiple(mocks)
			})

			It("finds the previous completed status and returns the current state with previous compelted state", func() {
				status, err := client.WorkflowStatus(pipelines, workflow)
				Ω(err).Should(BeNil())
				Ω(status).Should(Equal("canceled success"))
			})
		})
	})
})
