package cmd

import (
	"errors"
	"os"

	"github.com/YaleUniversity/deco/control"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [http(s)][ssm]://[some.host.name][/]path/to/control.json",
	Short: "Run executes the tasks in the given control file",
	Long: `Run executes the tasks passed in a control file.  Note the control file
can be a local file or an http/https endpoint and can be absolute or relative.  If no
control file is specified, the default '/var/run/secrets/deco.json' is assumed.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("accepts only one arg (the control location)")
		} else if len(args) == 1 {
			controlLocation = args[0]
		}
		// else we are using the default controlLocation
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var c control.Configuration
		if err := c.Read(controlLocation, httpHeaders, encoded); err != nil {
			Logger.Println("[ERROR] Unable to validate control file.", err)
			os.Exit(1)
		}

		// conditional here allows baseDir to be set in the JSON file
		// and not overriden with ""
		if baseDir != "" {
			c.BaseDir = baseDir
		}

		if err := c.DoFilters(); err != nil {
			Logger.Println("[ERROR] Filtering failed!", err)
			os.Exit(2)
		}
	},
	TraverseChildren: true,
}

func init() {
	runCmd.Flags().BoolVarP(&encoded, "encoded", "e", false, "Control file is base64 encoded")
	runCmd.Flags().StringArrayVarP(&httpHeaders, "header", "H", []string{}, "Pass a custom header to server")
	RootCmd.AddCommand(runCmd)
}
