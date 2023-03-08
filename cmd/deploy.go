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
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			panic(err)
		}
		deployWrapper(cmd, args, force)
	},
}
var deployCron = &cobra.Command{
	Use:   "cron",
	Short: "Continous Delivery with cron. It works by checking git repo status at interval",
	Run: func(cmd *cobra.Command, args []string) {
		c := cron.New()
		c.AddFunc("0 0 * * *", func() {
			deployWrapper(cmd, args, false)
		})
		c.Run()
	},
}

func deployWrapper(cmd *cobra.Command, args []string, force bool) {
	filepaths, err := cmd.Flags().GetStringArray("filepaths")
	if err != nil {
		log.Println(err)
		return
	}
	deploy(".", filepaths, force)

}

func deploy(workdir string, composePaths []string, force bool) {
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
	if strings.Contains(stringified, "branch is up to date") && !force {
		log.Println("No changes")
	} else {
		log.Println("Not in sync, pulling changes")
		gitPull := exec.Command("git", "pull")
		gitPull.Dir = workdir
		_, err := gitPull.Output()
		if err != nil {
			panic(err)
		}

		composeArgs := []string{
			"compose",
		}
		for _, s := range composePaths {
			composeArgs = append(composeArgs, "-f", s)
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

	}
}
