package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jerempy/brang/config"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "first time setup helper",
	Long: `
Puts data and returns response to terminal.
Accepts 1 positional arg of either a valid URL or a request saved in the .brang.yml using dot.notation.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if config.HomeDir == "" {
			fmt.Println("trouble finding home dir. exiting")
			return
		}
		fmt.Printf("This will create a 'brang' folder at %v, moves the executable to it, and creates the config files.\nThen it will offer to add brang to your path (so can be run in command line)\n", config.HomeDir)
		r := bufio.NewReader(os.Stdin)
		confirm(r, true)
		var bin string
		if runtime.GOOS == "windows" {
			bin = filepath.Join(config.BrangPath, "brang.exe")
		} else {
			bin = filepath.Join(config.BrangPath, "brang")
		}
		os.MkdirAll(config.ConfigPath, os.ModePerm)
		fmt.Printf("Add new config files. This will overwrite any that already exist at %v. Say 'y' if doing first time setup.\n", config.BrangPath)
		skip := confirm(r, false)
		if !skip {
			err := os.WriteFile(config.ConfigFile, configTmpl, os.ModePerm)
			if err != nil {
				fmt.Println("err creating config file", err)
				return
			}
			err = os.WriteFile(config.RequestsFile, requestsTmpl, os.ModePerm)
			if err != nil {
				fmt.Println("err creating requests file", err)
				return
			}
			fmt.Println("Added config files")
		}
		// move binary to BrangPath
		cur, err := os.Executable()
		if err != nil {
			panic(err)
		}
		os.Rename(cur, bin)

		// Adds dir to path so brang can be run in shell
		fmt.Println("Allow brang to try to set up your path variable so you can execute brang in the shell?")
		skip = confirm(r, false)
		if !skip {
			if runtime.GOOS == "windows" {
				ps, err := exec.LookPath("powershell.exe")
				if err != nil {
					fmt.Println(err)
					return
				}
				psCmd := fmt.Sprintf(`[Environment]::SetEnvironmentVariable(
				"Path",
				[Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) + ";%v\brang",
				[EnvironmentVariableTarget]::User)`, config.HomeDir)
				setPath := exec.Command(ps, psCmd)
				setPath.Run()
			} else {
				setPath := exec.Command("ln", "-s", bin, "/usr/local/bin/brang")
				setPath.Run()
				fmt.Printf("The binary is at %v, there is a symlink at /usr/local/bin. If usr/local/bin is not in your path you'll need to add it. Google for explanation for your system.", bin)
			}
		}
		fmt.Println("Setup all done! You might need to restart any terminal/shell to start using. If its not working might need to ensure brang is in your path environment variable.")
		time.Sleep(8 * time.Second)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func confirm(r *bufio.Reader, exit bool) (skipNext bool) {
	fmt.Print("'y' to confirm, anything else not to: ")
	t, err := r.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	text := string(t[0])
	if text != "y" && text != "Y" {
		if exit {
			fmt.Println("didnt confirm. Exiting")
			os.Exit(0)
		} else {
			fmt.Println("didnt confirm. skipping.")
			skipNext = true
		}
	}
	return
}

var configTmpl = []byte(`# Brang configuration options
outWriter: stdout # stdout|file|tempFile
outWriterFormat: pretty # pretty|basic|raw
deleteTempFileOnClose: true
# outWriterFileType: txt #full named path of file
# outWriterFileName: brangoutput
# outWriterFilePath: /usr
# fileEditor: notepad #Either alias or direct path to executable
`)

var requestsTmpl = []byte(`# Saved requests
# These requests are then accessed in cmd line like: mysite.posts.all or github.brangreadme
# Anything can be a reference to a env variable like: $THE_VAR - just make sure its set in your environment.
# Examples:
# mysite:
#   auth:
#     authtype: Bearer
#     token: ABC-456 # could use $MYSITE_TOKEN
#     username: joe # could use $MYSITE_USERNAME
#     password: secret # could use $MYSITE_PASSWORD
#   requests:
#     users: https://mysite.com/users/
#     posts:
#       all:
#         url: https://mysite.com/posts/
#         body: |
#           {
#             "test": "test", 
#             "body": "this here is a big ol body.",
#             "userId": 27
#           }
#
#         header:
#           accept: text/plain,text/html
#           content-type: application/json

#       firstpost: https://jsonplaceholder.typicode.com/posts/1

# github:
#   requests:
#     brangreadme: https://raw.githubusercontent.com/jerempy/brang/main/README.md
`)
