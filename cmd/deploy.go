package cmd

import (
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
		deployWrapper(cmd, args)
	},
}
var deployCron = &cobra.Command{
	Use:   "cron",
	Short: "Continous Delivery with cron. It works by checking git repo status at interval",
	Run: func(cmd *cobra.Command, args []string) {
		c := cron.New()
		c.AddFunc("0 0 * * *", func() {
			deployWrapper(cmd, args)
		})
		c.Run()
	},
}

func deployWrapper(cmd *cobra.Command, args []string) {
	filepaths, err := cmd.Flags().GetStringArray("filepaths")
	if err != nil {
		log.Println(err)
		return
	}
	deploy(".", filepaths)

}

func deploy(workdir string, composePaths []string) {
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
		log.Println("No changes")
	} else {
		log.Println("Not in sync, pulling changes")
		gitPull := exec.Command("git", "pull")
		gitPull.Dir = workdir
		gitPullOut, err := gitPull.Output()
		if err != nil {
			panic(err)
		}
		log.Println(string(gitPullOut))

		composePathsString := ""
		for _, s := range composePaths {
			composePathsString += "-f " + s + " "
		}
		composePathsString = strings.TrimSpace(composePathsString)
		composeUp := exec.Command("docker", "compose", composePathsString, "up", "-d", "--build", "--remove-orphans")
		composeUp.Stdout = os.Stdout
		log.Println("Args: ", composeUp.Args)
		composeUp.Dir = workdir
		err = composeUp.Start()
		if err != nil {
			panic(err)
		}

	}
}
