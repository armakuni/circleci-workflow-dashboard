package circleci_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
)

var _ = Describe("Projects", func() {
	var project = circleci.Project{
		VCSType:  "github",
		Username: "foobar",
		Reponame: "example",
	}

	var projects = circleci.Projects{
		{
			VCSType:  "github",
			Username: "foobar",
			Reponame: "example",
		},
		{
			VCSType:  "github",
			Username: "another",
			Reponame: "example",
		},
		{
			VCSType:  "github",
			Username: "foobar",
			Reponame: "another",
		},
		{
			VCSType:  "github",
			Username: "fwibble",
			Reponame: "fish",
		},
	}
	var filter = &circleci.Filter{
		"foobar/example": nil,
		"fwibble/fish":   nil,
	}

	Describe("#Slug", func() {
		It("returns the circleci project slug", func() {
			Ω(project.Slug()).Should(Equal("github/foobar/example"))
		})
	})

	Describe("#Name", func() {
		It("returns the circleci project name", func() {
			Ω(project.Name()).Should(Equal("foobar/example"))
		})
	})

	Describe("#Filter", func() {
		Context("if the filter is blank", func() {
			It("returns all prokects", func() {
				Ω(projects).Should(HaveLen(4))
				filteredProjects := projects.Filter(&circleci.Filter{})
				Ω(filteredProjects).Should(HaveLen(4))
			})
		})

		Context("if the filter is populated", func() {
			It("returns only the projects with matching names", func() {
				Ω(projects).Should(HaveLen(4))
				filteredProjects := projects.Filter(filter)
				Ω(filteredProjects).Should(HaveLen(2))
			})
		})
	})
})
