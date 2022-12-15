package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/YaleUniversity/deco/control"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [http(s)][ssm]://[some.host.name][/]path/to/control.json",
	Short: "Validates the control file",
	Long: `Validates the control file format.  Note the control file
can be a local file or an http/https endpoint and can be absolute or relative. If no
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
		if err := c.Read(controlLocation, httpHeaders, httpTimeout, encoded); err != nil {
			Logger.Println("[ERROR] Unable to validate control file.", err)
			os.Exit(1)
		}
		Logger.Println("Control file validated successfully.")
	},
	TraverseChildren: true,
}

func init() {
	validateCmd.Flags().BoolVarP(&encoded, "encoded", "e", false, "Control file is base64 encoded")
	validateCmd.Flags().StringArrayVarP(&httpHeaders, "header", "H", []string{}, "Pass a custom header to server")
	validateCmd.Flags().DurationVarP(&httpTimeout, "timeout", "T", 15*time.Second, "Set the HTTP request timeout")
	RootCmd.AddCommand(validateCmd)
}
