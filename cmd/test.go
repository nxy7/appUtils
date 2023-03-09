package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
	Short: `Runs compose files provided by user with --testfile | -t flags and watches the output of services containing 'test'. If any
		test container fails.`,
	Run: func(cmd *cobra.Command, args []string) {
		res := test(cmd)
		if !res {
			os.Exit(1)
		}
	},
}

// Returns true if tests are completed successfully
func test(cmd *cobra.Command) bool {
	// run docker compose with all necesarry compose files

	// watch output for all containers with "test in name"

	// when all containers end their jobs collect results

	// if any single one fails return false

	// else return true
	return true
}
