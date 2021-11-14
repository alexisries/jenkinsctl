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

package job

import (
	"bufio"
	"errors"
	"fmt"
	"jenkinsctl/pkg/apiclient"
	"jenkinsctl/pkg/apiclient/jobs"
	"os"

	"github.com/spf13/cobra"
)

func NewJobCmd(client *apiclient.ApiClient) *cobra.Command {

	// cmd represents the job command
	var cmd = &cobra.Command{
		Use:   "job",
		Short: "This command allows to manage a job on Jenkins",
		Long: `This command allows to manage a job on Jenkins

For example:

list jobs:
	jenkinsctl job list
	jenkinsctl job list --minimum-age=1h
	jenkinsctl job list --name=my-app

start jobs:
	jenkinsctl job start --name=my-app
	jenkinsctl job start --minimum-age=1h
	jenkinsctl job start --name=my-app --schedule="0/2/*/*/*"

stop jobs:
	jenkinsctl job stop --minimum-age=1h
	jenkinsctl job stop --name=my-app`,
	}

	cmd.AddCommand(NewJobListCmd(client))
	cmd.AddCommand(NewJobStartCmd(client))
	cmd.AddCommand(NewJobStopCmd(client))
	return cmd
}

func checkStatusValidValue(status string) error {
	if status != jobs.JOB_STATUS_ALL &&
		status != jobs.JOB_STATUS_RUNNING &&
		status != jobs.JOB_STATUS_SUCCESS &&
		status != jobs.JOB_STATUS_FAILED &&
		status != jobs.JOB_STATUS_ABORTED {
		return fmt.Errorf("%s is not accepted status", status)
	}
	return nil
}

func askUserForYesOrNo(action string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nDo you want to %s these jobs ? (yes or no): ", action)
	userInput, _ := reader.ReadString('\n')
	if userInput == "no\n" {
		return errors.New("user canceled")
	} else if userInput != "yes\n" {
		return fmt.Errorf("unrecognized command: %s", userInput)
	}
	fmt.Println()
	return nil
}
