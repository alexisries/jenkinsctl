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
	"fmt"
	"jenkinsctl/pkg/apiclient"
	"jenkinsctl/pkg/apiclient/jobs"

	"github.com/spf13/cobra"
)

type JobStopFlags struct {
	Name      string
	AgeMin    int
	AgeMax    int
	ForceStop bool
}

func newJobStopFlags() *JobStopFlags {
	return &JobStopFlags{
		Name:      "",
		AgeMin:    0,
		AgeMax:    0,
		ForceStop: false,
	}
}

func NewJobStopCmd(client *apiclient.ApiClient) *cobra.Command {

	jobStopFlags := newJobStopFlags()

	// cmd represents the job command
	var cmd = &cobra.Command{
		Use:   "stop",
		Short: "stop jobs",
		Long: `For example:
	jenkinsctl job stop --minimum-age=1h
	jenkinsctl job stop --name=my-app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return jobStop(client, jobStopFlags)
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(
		&jobStopFlags.Name, "name", jobStopFlags.Name,
		"Filter Jobs from the name",
	)
	cmd.Flags().IntVar(
		&jobStopFlags.AgeMin, "minimum-age", jobStopFlags.AgeMin,
		"Filter Jobs from last build minimum age (in minutes)",
	)
	cmd.Flags().IntVar(
		&jobStopFlags.AgeMax, "maximum-age", jobStopFlags.AgeMax,
		"Filter Jobs from last build maximum age (in minutes)",
	)
	cmd.Flags().BoolVar(
		&jobStopFlags.ForceStop, "force", jobStopFlags.ForceStop,
		"Force stop jobs",
	)
	return cmd
}

func jobStop(client *apiclient.ApiClient, flags *JobStopFlags) error {
	filter := jobs.JobsFilterParams{
		Name:   flags.Name,
		AgeMin: flags.AgeMin,
		AgeMax: flags.AgeMax,
		Status: jobs.JOB_STATUS_RUNNING,
	}
	jobs := jobs.Jobs{}
	err := jobs.GetFilteredJobs(client, &filter)
	if err != nil {
		return err
	}

	if len(jobs.Jobs) == 0 {
		fmt.Println("all jobs are in stopped state")
		return nil
	}

	fmt.Println("\nJobs to be stopped :")
	jobs.PrintJobsTable()
	if !flags.ForceStop {
		err = askUserForYesOrNo("stop")
		if err != nil {
			return err
		}
	}
	return jobs.Stop(client, flags.ForceStop)
}
