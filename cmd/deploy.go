package cmd

import (
	"log"
	"os/exec"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Command used to deploy applications using docker compose files",
	Long:  ``,
}

var deployNow = &cobra.Command{
	Use:   "now",
	Short: "Check if there are any changes to repo and redeploy the app",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		filepaths, err := cmd.Flags().GetStringArray("filepaths")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(filepaths)
		log.Println("deploy now")
		deploy("..", filepaths)
	},
}
var deployCron = &cobra.Command{
	Use:   "cron",
	Short: "Continous Delivery with cron, that works by checking git repo status at interval",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		c := cron.New()
		c.AddFunc("0 0 * * *", func() {
			// deploy("..")
		})
		c.Run()
	},
}

func deployWrapper(cmd *cobra.Command, args []string)

func deploy(workdir string, composePaths []string) {
	lsCmd := exec.Command("ls")
	lsCmd.Dir = workdir
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	log.Println(string(lsOut))
	gitCmd := exec.Command("git", "status")
	gitOut, err := gitCmd.Output()
	if err != nil {
		panic(err)
	}
	log.Println(string(gitOut))
}
