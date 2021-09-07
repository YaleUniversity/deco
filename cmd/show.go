package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/YaleUniversity/deco/control"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [http(s)][ssm]://[some.host.name][/]path/to/control.json",
	Short: "Reads and displays a control file on STDOUT",
	Long: `Show reads and displays a control file on STDOUT.  Note the control file
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
		if err := c.Read(controlLocation, httpHeaders, encoded); err != nil {
			Logger.Println("[ERROR] Unable to read control file.", err)
			os.Exit(1)
		}

		b, err := json.MarshalIndent(&c, "", "    ")
		if err != nil {
			Logger.Println("[ERROR] Unable to marshal control file.", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", string(b))
	},
	TraverseChildren: true,
}

func init() {
	showCmd.Flags().BoolVarP(&encoded, "encoded", "e", false, "Control file is base64 encoded")
	showCmd.Flags().StringArrayVarP(&httpHeaders, "header", "H", []string{}, "Pass a custom header to server")
	RootCmd.AddCommand(showCmd)
}
