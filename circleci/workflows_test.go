package circleci_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/armakuni/circleci-workflow-dashboard/circleci"
)

var _ = Describe("Workflow", func() {
	Describe("#BuildError", func() {
		Context("when the workflow had a build error", func() {
			It("returns true", func() {
				workflow := circleci.Workflow{Name: "Build Error"}
				Ω(workflow.BuildError()).Should(BeTrue())
			})
		})

		Context("when the workflow did not have a build error", func() {
			It("returns false", func() {
				workflow := circleci.Workflow{Name: "foobar"}
				Ω(workflow.BuildError()).Should(BeFalse())
			})
		})
	})
})

var _ = Describe("Workflows", func() {
	Describe("#BuildError", func() {
		Context("when at least 1 workflow had a build error", func() {
			It("returns true", func() {
				workflows := circleci.Workflows{{Name: "Build Error"}}
				Ω(workflows.BuildError()).Should(BeTrue())
			})
		})

		Context("when no workflows did not have a build error", func() {
			It("returns false", func() {
				workflows := circleci.Workflows{{Name: "foobar"}}
				Ω(workflows.BuildError()).Should(BeFalse())
			})
		})
	})
})
