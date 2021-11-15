/*
Copyright © 2021 Alexis Ries <ries.alexis@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"jenkinsctl/pkg/apiclient"
	"jenkinsctl/pkg/cmd/job"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(rootCmd *cobra.Command) {
	cobra.CheckErr(rootCmd.Execute())
}

func NewRootCmd(client *apiclient.ApiClient) *cobra.Command {

	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(client.Initialize)

	// cmd represents the base command when called without any subcommands
	var cmd = &cobra.Command{
		Use:   "jenkinsctl",
		Short: "CLI for managing jobs on Jenkins ",
		Long: `With this command-line tool you will be able to list, start and stop jobs on Jenkins.
This tool was developed as part of a technical challenge for Bonitasoft.`,
	}

	cmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.jenkinsctl.yaml)",
	)

	cmd.AddCommand(job.NewJobCmd(client))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	return cmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		//Find current directory
		currentDirectory, err := os.Getwd()
		cobra.CheckErr(err)

		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".jenkinsctl" (without extension).
		viper.AddConfigPath(currentDirectory)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jenkinsctl")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
