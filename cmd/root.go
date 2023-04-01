package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/inconshreveable/mousetrap"
	"github.com/jerempy/brang/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "brang",
	Short: "A CLI tool for HTTP requests",
	Long: `brang - like instead of saying brought:

This is a CLI tool for simplifying HTTP requests.
It brings additional functionality to requests such as
saving credentials and saving HTTP requests for repeat use.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("You should run brang from command line - unless this is first time setup.")
			fmt.Println("If first time setup - then confirm to start setup installer")
			setupCmd.Execute()
		}
	}
	config.Requests.SetConfigFile(config.RequestsFile)
	err := config.LoadBrangConfig()
	if err != nil {
		fmt.Printf("couldn't find config files. err: %v -- Confirm to start setup installer\n", err)
		setupCmd.Execute()
		return
	}
}
