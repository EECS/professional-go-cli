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

var configFileName string

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
		// Get managed repos from viper
		root := viper.GetString("root")
		if root[len(root)-1] != '/' {
			root += "/"
		}

		repoNames := []string{}
		if len(viper.GetString("repos")) > 0 {
			repoNames = strings.Split(viper.GetString("repos"), ",")
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

func initConfig() {
	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		panic(err)
	}

	viper.SetConfigFile(configFileName)
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("multi_git")
	err = viper.BindEnv("root")
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv("repos")
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv("ignore-errors")
	if err != nil {
		panic(err)
	}
}

func init() {
	parentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defaultConfigFileName := path.Join(parentDir, "/multi-git.toml")

	rootCmd.Flags().StringVar(&configFileName,
		"config", defaultConfigFileName,
		"config file path (default is parentDir/multi-git.toml)")

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().Bool(
		"ignore-errors",
		false,
		`will continue executing the command for all repos if ignore-errors
                 is true otherwise it will stop execution when an error occurs`)

	err = viper.BindPFlag("ignore-errors", rootCmd.Flags().Lookup("ignore-errors"))
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
