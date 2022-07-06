package repo_manager_test

import (
	"os"
	"path"
	"strings"
	"testing"

	. "multi-git/pkg/helpers"

	. "multi-git/pkg/repo_manager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const baseDir = "tmp/test-multi-git"

var repoList = []string{}

func TestRepoManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RepoManager Suite")
}

var _ = Describe("Repo manager tests", func() {
	var err error

	removeAll := func() {
		err = os.RemoveAll(baseDir)
		Expect(err).Should(BeNil())
	}

	BeforeEach(func() {
		removeAll()
		err = CreateDir(baseDir, "dir-1", true)
		Expect(err).Should(BeNil())
		repoList = []string{"dir-1"}
	})
	AfterEach(removeAll)

	Context("Tests for failure cases", func() {
		It("Should fail with invalid base dir", func() {
			_, err := NewRepoManager("/no-such-dir", repoList, true)
			Expect(err).ShouldNot(BeNil())
		})
		It("Should fail with empty repo list", func() {
			_, err := NewRepoManager(baseDir, []string{}, true)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("Tests for success cases", func() {
		It("Should commit files successfully", func() {
			rm, err := NewRepoManager(baseDir, repoList, true)
			Expect(err).Should(BeNil())

			output, err := rm.Exec("checkout -b test-branch")
			Expect(err).Should(BeNil())

			for _, out := range output {
				Expect(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
			}

			AddFiles(baseDir, repoList[0], true, "file_1.txt", "file_2.txt")

			// Restore working directory after executing the command
			wd, _ := os.Getwd()
			defer os.Chdir(wd)

			dir := path.Join(baseDir, repoList[0])
			err = os.Chdir(dir)
			Expect(err).Should(BeNil())

			output, err = rm.Exec("log --oneline")
			Expect(err).Should(BeNil())

			ok := strings.HasSuffix(output[dir], "added some files...\n")
			Expect(ok).Should(BeTrue())
		})
	})
})
