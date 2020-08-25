package circleci_test

import (
	"github.com/armakuni/circleci-workflow-dashboard/circleci"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipelines", func() {
	var pipelines = circleci.Pipelines{
		{
			ID:  "1",
			VCS: circleci.VCS{Branch: "master"},
		},
		{
			ID:  "2",
			VCS: circleci.VCS{Branch: "master"},
		},
		{
			ID:  "3",
			VCS: circleci.VCS{Branch: "master"},
		},
		{
			ID:  "4",
			VCS: circleci.VCS{Branch: "develop"},
		},
		{
			ID:  "5",
			VCS: circleci.VCS{Branch: "develop"},
		},
	}

	Describe("#FilterPerBranch", func() {
		Context("when you set a branch filter", func() {
			It("returns pipelines in a map keyed by branch", func() {
				filteredPipelines := pipelines.FilteredPerBranch("master")
				Ω(filteredPipelines).Should(HaveLen(1))
				Ω(filteredPipelines["master"]).Should(HaveLen(3))
			})
		})

		Context("when you don't set a branch filter", func() {
			It("returns pipelines in a map keyed by branch", func() {
				filteredPipelines := pipelines.FilteredPerBranch("")
				Ω(filteredPipelines).Should(HaveLen(2))
				Ω(filteredPipelines["master"]).Should(HaveLen(3))
				Ω(filteredPipelines["develop"]).Should(HaveLen(2))
			})
		})
	})

	Describe("#LatestPerBranch", func() {
		It("returns the latest pipelines for each branch", func() {
			filteredPipelines := pipelines.LatestPerBranch()
			Ω(filteredPipelines).Should(HaveLen(2))
			Ω(filteredPipelines["master"].ID).Should(Equal("1"))
			Ω(filteredPipelines["develop"].ID).Should(Equal("4"))
		})
	})
})
