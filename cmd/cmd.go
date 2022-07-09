package cmd

import (
	"fmt"
	"log"
	"multi-git/pkg/repo_manager"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "multi-git",
	Short: "Runs git commands over multiple repos",
	Long: `Runs git commands over multiple repos.

Requires the following environment variables defined:   
MG_ROOT: root directory of target git repositories
MG_REPOS: list of repository names to operate on`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get managed repos from environment variables
		root := os.Getenv("MG_ROOT")
		if root[len(root)-1] != '/' {
			root += "/"
		}

		repoNames := []string{}
		if len(os.Getenv("MG_REPOS")) > 0 {
			repoNames = strings.Split(os.Getenv("MG_REPOS"), ",")
		}

		repoManager, err := repo_manager.NewRepoManager(root, repoNames, viper.GetBool("ignore-errors"))
		if err != nil {
			log.Fatal(err)
		}

		command := strings.Join(args, " ")
		output, err := repoManager.Exec(command)
		if err != nil {
			fmt.Printf("command '%s' failed with error ", err)
		}

		for repo, out := range output {
			fmt.Printf("[%s]: git %s\n", path.Base(repo), command)
			fmt.Println(out)
		}
	},
}

func init() {
	rootCmd.Flags().Bool(
		"ignore-errors",
		false,
		`will continue executing the command for all repos if ignore-errors
                 is true otherwise it will stop execution when an error occurs`)

	err := viper.BindPFlag("ignore-errors", rootCmd.Flags().Lookup("ignore-errors"))
	if err != nil {
		panic("Unable to bind flag")
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
