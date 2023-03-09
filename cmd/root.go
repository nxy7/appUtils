/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "apputils",
	Short: "Utility 'scripts' helpful for range of applications",
	Long: `Utility 'scripts' made to automate application deployment on server
		and creating cicd pipelines. It requires git, docker and 'docker compose' commands to be
		available on the system`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(deployCmd)
	// deployCmd.PersistentFlags().StringP("workdir", "w", ".", "Path to use as working directory")
	deployCmd.PersistentFlags().StringArrayP("filepaths", "f", []string{}, "Paths to docker compose files used for production build.")
	deployCmd.PersistentFlags().StringP("workdir", "w", ".", "Workdir that the command should be ran in.")
	deployCmd.PersistentFlags().StringArrayP("testfile", "t", []string{}, `Paths to docker compose files used for testing. If present,
		apputils will first test against those files. Using 'name' propety in compose files is highly recommended, so testing doesn't interfere
		with current app deployments.`)
	deployCmd.PersistentFlags().String("project-name", "", `Project name to use with docker compose command. If supplied and testfiles are present
		apputils will use 'project-name' postfixed with _test as project name for tests, to avoid clashing names.`)
	deployCmd.MarkFlagRequired("filepaths")
	deployCmd.AddCommand(deployNow, deployCron, testCmd)
	deployCron.PersistentFlags().StringP("cron", "c", "30 3 * * *", "Deploy interval in cron format, defaults to 30 3 * * *")
	deployNow.Flags().Bool("force", false, "Force Redeployment of the app")
}
