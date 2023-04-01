package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/jerempy/brang/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "view or edit configurations",
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your brang configurations",
	Run: func(cmd *cobra.Command, args []string) {
		f := config.Brang.ConfigFileUsed()
		d, err := os.ReadFile(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(d))
	},
}

var configWhereCmd = &cobra.Command{
	Use:   "where",
	Short: "List locations of brang configurations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("brang path - where everything is located: ", config.BrangPath)
		fmt.Println("config.yaml - main settings: ", config.ConfigFile)
		fmt.Println("requests.yaml - saved requests: ", config.RequestsFile)
	},
}

var configOpenCmd = &cobra.Command{
	Use:       "open [path] [config] [requests]",
	Short:     "Open the config or path for viewing/editing",
	ValidArgs: []string{"path", "config", "requests"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var p string
		switch args[0] {
		case "path":
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("explorer", config.BrangPath)
			} else if runtime.GOOS == "darwin" {
				cmd = exec.Command("open", config.BrangPath)
			} else {
				fmt.Printf("type: 'cd %v'", config.BrangPath)
				return
			}
			cmd.Run()
			return
		case "config":
			p = config.ConfigFile
		case "requests":
			p = config.RequestsFile
		default:
			fmt.Println("Couldn't find file. use --help (-h) flag for choices")
			return
		}
		config.OpenEditor(p).Run()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd, configWhereCmd, configOpenCmd)
	configOpenCmd.Flags().StringP("editor", "e", "", "either alias name or path to executable for your editor")
	config.Brang.BindPFlag("fileEditor", configOpenCmd.Flags().Lookup("editor"))
}
