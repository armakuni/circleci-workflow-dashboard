package circleci_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCircleCI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CircleCI Suite")
}
