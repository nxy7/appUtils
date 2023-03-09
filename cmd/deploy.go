package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Command used to deploy applications using docker compose files",
}

var deployNow = &cobra.Command{
	Use:   "now",
	Short: "Check if there are any changes to repo and redeploy the app",
	Run: func(cmd *cobra.Command, args []string) {
		deploy(cmd)
	},
}
var deployCron = &cobra.Command{
	Use:   "cron",
	Short: "Continous Delivery with cron. It works by checking git repo status at interval",
	Run: func(cmd *cobra.Command, args []string) {
		c := cron.New()
		c.AddFunc("0 0 * * *", func() {
			deploy(cmd)
		})
		c.Run()
	},
}

func deploy(cmd *cobra.Command) error {
	composePaths, err := cmd.Flags().GetStringArray("filepaths")
	if err != nil {
		panic(err)
	}
	testPaths, err := cmd.Flags().GetStringArray("testfile")
	if err != nil {
		panic(err)
	}
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		panic(err)
	}
	projectName, err := cmd.Flags().GetString("project-name")
	if err != nil {
		panic(err)
	}
	workdir, err := cmd.Flags().GetString("project-name")
	if err != nil {
		panic(err)
	}

	if isBranchUpToDate(workdir) && !force {
		log.Println("No changes")
		return nil
	}

	log.Println("Not in sync, pulling changes")
	gitPullChanges(workdir)

	if len(testPaths) > 0 {
		if !test(cmd) {
			return fmt.Errorf("App didn't pass tests")
		}
	}

	composeArgs := []string{
		"compose",
	}
	for _, s := range composePaths {
		composeArgs = append(composeArgs, "-f", s)
	}
	if projectName != "" {
		composeArgs = append(composeArgs, "--project-name", projectName)
	}
	composeArgs = append(composeArgs, "up", "-d", "--build", "--remove-orphans")
	if force {
		composeArgs = append(composeArgs, "--force-recreate")
	}

	composeUp := exec.Command("docker", composeArgs...)
	composeUp.Stdout = os.Stdout
	log.Println("Args: ", composeUp.Args)
	composeUp.Dir = workdir
	err = composeUp.Run()
	if err != nil {
		panic(err)
	} else {
		log.Println("App deployed")
	}
	return nil
}

func isBranchUpToDate(workdir string) bool {

	gitFetch := exec.Command("git", "fetch")
	gitFetch.Dir = workdir
	err := gitFetch.Run()
	if err != nil {
		panic(err)
	}

	gitCmd := exec.Command("git", "status")
	gitCmd.Dir = workdir
	gitOut, err := gitCmd.Output()
	if err != nil {
		panic(err)
	}
	stringified := string(gitOut)

	if strings.Contains(stringified, "branch is up to date") {
		return true
	}
	return false
}

func gitPullChanges(workdir string) {
	gitPull := exec.Command("git", "pull")
	gitPull.Dir = workdir
	_, err := gitPull.Output()
	if err != nil {
		panic(err)
	}
}
