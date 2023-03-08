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
	deployCmd.PersistentFlags().StringArrayP("filepaths", "f", []string{}, "paths to docker compose files used for production build")
	deployCmd.MarkFlagRequired("filepaths")
	deployCmd.AddCommand(deployNow, deployCron)
	deployNow.Flags().Bool("force", false, "Force Redeployment of the app")
}
