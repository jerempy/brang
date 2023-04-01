package cmd

import (
	"fmt"
	"strings"

	"github.com/jerempy/brang/client"
	"github.com/spf13/cobra"
)

var rset client.RequestSet

var getCmd = &cobra.Command{
	Use:   "get {url|SavedRequest}",
	Short: "HTTP POST Request",
	Long: `
Gets the data and returns to terminal.
Accepts 1 positional arg of either a valid URL or a request saved in the .brang.yml using dot.notation.`,
	Example: `'brang get https://mysite.com/users' or using SavedRequests: 'brang get mysite.users'`,
	Args:    cobra.ExactArgs(1),
	Run:     processAndRunRequest,
}

var postCmd = &cobra.Command{
	Use:   "post {url|SavedRequest}",
	Short: "HTTP POST Request",
	Long: `
Posts data and returns response to terminal.
Accepts 1 positional arg of either a valid URL or a request saved in the .brang.yml using dot.notation.`,
	Example: `'brang post https://mysite.com/users' or using SavedRequests: 'brang post mysite.users'`,
	Args:    cobra.ExactArgs(1),
	Run:     processAndRunRequest,
}

var putCmd = &cobra.Command{
	Use:   "put {url|SavedRequest}",
	Short: "HTTP PUT Request",
	Long: `
Puts data and returns response to terminal.
Accepts 1 positional arg of either a valid URL or a request saved in the .brang.yml using dot.notation.`,
	Example: `'brang put https://mysite.com/users' or using SavedRequests: 'brang put mysite.users'`,
	Args:    cobra.ExactArgs(1),
	Run:     processAndRunRequest,
}

var patchCmd = &cobra.Command{
	Use:   "patch {url|SavedRequest}",
	Short: "HTTP PATCH Request",
	Long: `
Patch data and returns response to terminal.
Accepts 1 positional arg of either a valid URL or a request saved in the .brang.yml using dot.notation.`,
	Example: `'brang patch https://mysite.com/users' or using SavedRequests: 'brang patch mysite.users'`,
	Args:    cobra.ExactArgs(1),
	Run:     processAndRunRequest,
}

var deleteCmd = &cobra.Command{
	Use:   "delete {url|SavedRequest}",
	Short: "HTTP DELETE Request",
	Long: `
Delete data and returns response to terminal.
Accepts 1 positional arg of either a valid URL or a request saved in the .brang.yml using dot.notation.`,
	Example: `'brang delete https://mysite.com/users' or using SavedRequests: 'brang delete mysite.users'`,
	Args:    cobra.ExactArgs(1),
	Run:     processAndRunRequest,
}

func init() {
	rootCmd.AddCommand(getCmd, putCmd, postCmd, patchCmd, deleteCmd)
	requestCmdFlags(getCmd)
	requestCmdFlags(putCmd)
	requestCmdFlags(postCmd)
	requestCmdFlags(patchCmd)
	requestCmdFlags(deleteCmd)
}

func requestCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&rset.AuthType, "auth", "a", "", "set Auth type: Password|Token|Bearer")
	cmd.Flags().StringVarP(&rset.Cred, "cred", "c", "", `set token for auth types Token|Bearer ex: 123-456-ABC.
or set username:password for type Password ex: john123:secretpass`)
	cmd.Flags().StringArrayVarP(&rset.HeaderSlice, "header", "H", []string{}, `set headers as key:value, as many as needed. 
ex: -H "Content-Type:application/json" -H "Authorization:Bearer 123-456-ABC"`)
	cmd.Flags().StringVarP(&rset.Params, "params", "p", "", `attaches additional params to url. ex for https://mysite.com:
-p ?title=MyTitle -> https://mysite.com/?title=Mytitle or -p 123 ->https://mysite.com/123`)
	cmd.Flags().StringVarP(&rset.Body, "body", "b", "", `includes body in the request, such as json body for a post request.`)
	cmd.Flags().StringP("file", "f", "", `path to a file for body of request`)
}

func processAndRunRequest(cmd *cobra.Command, args []string) {
	rset.Method = strings.ToUpper(cmd.Name())
	rset.URL = args[0]
	if file, _ := cmd.Flags().GetString("file"); file != "" {
		if err := rset.BodyFile(file); err != nil {
			fmt.Printf("err reading body file: %v", err)
			return
		}
	}
	rset.Send()
}
