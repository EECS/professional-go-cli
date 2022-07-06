package e2e_tests

import (
	"os"

	. "multi-git/pkg/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const baseDir = "/tmp/multi-git"

var repoList string

var _ = Describe("multi-git e2e tests", func() {
	var err error

	removeAll := func() {
		err = os.RemoveAll(baseDir)
		Expect(err).Should(BeNil())
	}

	BeforeEach(func() {
		removeAll()
		err = CreateDir(baseDir, "", false)
		Expect(err).Should(BeNil())
	})

	AfterSuite(removeAll)
})
