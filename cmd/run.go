// Copyright © 2017 Yale University
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"os"

	"github.com/YaleUniversity/deco/control"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run executes the taks in the given control file",
	Long:  "",
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
		if err := c.Read(controlLocation); err != nil {
			Logger.Println("Unable to validate control file", err)
			os.Exit(1)
		}

		// conditional here allows baseDir to be set in the JSON file
		// and not overriden with ""
		if baseDir != "" {
			c.BaseDir = baseDir
		}

		if err := c.DoFilters(); err != nil {
			Logger.Println("Filtering failed!", err)
			os.Exit(2)
		}
	},
}

var Source string

func init() {
	runCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
	RootCmd.AddCommand(runCmd)
}
