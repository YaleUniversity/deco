// Copyright © 2020 Yale University
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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	controlLocation = "/var/run/secrets/deco.json"
	cfgFile         string
	baseDir         string
	httpHeaders     []string
	encoded         bool
)

// Logger is a STDERR logger
var Logger = log.New(os.Stderr, "", 0)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "deco",
	Short: "Gets stuff ready to run inside docker.",
	Long: `deco gets your app ready to run when a container
starts.  For example: the filters allow you to specify
individual files to filter and key/value pairs to use when
filtering.  By default, it works from the current directory and
will filter files in place.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		Logger.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "deco config file -- _not_ the control file (default is $HOME/.deco.yaml)")
	RootCmd.PersistentFlags().StringVarP(&baseDir, "dir", "d", "", "Base directory for filtered files/templates")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".deco")           // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME")) // adding home directory as first search path

	viper.SetEnvPrefix("deco") // prefix environment variables with DECO_
	viper.AutomaticEnv()       // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		Logger.Println("Using config file:", viper.ConfigFileUsed())
	}
}
