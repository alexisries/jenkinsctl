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
	"errors"
	"fmt"
	"jenkinsctl/pkg/apiclient"
	"jenkinsctl/pkg/apiclient/jobs"

	"github.com/spf13/cobra"
)

type JobStartFlags struct {
	Name       string
	Cron       string
	AgeMin     int
	AgeMax     int
	Status     string
	ForceStart bool
}

func newJobStartFlags() *JobStartFlags {
	return &JobStartFlags{
		Name:       "",
		AgeMin:     0,
		AgeMax:     0,
		Status:     "all",
		Cron:       "",
		ForceStart: false,
	}
}

func NewJobStartCmd(client *apiclient.ApiClient) *cobra.Command {
	jobStartFlags := newJobStartFlags()

	// cmd represents the job command
	var cmd = &cobra.Command{
		Use:   "start",
		Short: "start a job",
		Long: `For example:
	jenkinsctl job start --name=my-app
	jenkinsctl job start --minimum-age=1h
	jenkinsctl job start --name=my-app --schedule=@daily`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return jobStart(client, jobStartFlags)
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(
		&jobStartFlags.Name, "name", jobStartFlags.Name,
		"Filter Jobs from the name",
	)
	cmd.Flags().StringVar(
		&jobStartFlags.Status, "status", jobStartFlags.Status,
		"Filter Jobs from status (possible values: all, running, success, aborted, failure)",
	)
	cmd.Flags().IntVar(
		&jobStartFlags.AgeMin, "minimum-age", jobStartFlags.AgeMin,
		"Filter Jobs from last build minimum age (in minutes)",
	)
	cmd.Flags().IntVar(
		&jobStartFlags.AgeMin, "maximum-age", jobStartFlags.AgeMax,
		"Filter Jobs from last build maximum age (in minutes)",
	)
	cmd.Flags().StringVar(
		&jobStartFlags.Cron, "schedule", jobStartFlags.Cron,
		"Specify the schedule in Jenkins time trigger syntax",
	)
	cmd.Flags().BoolVar(
		&jobStartFlags.ForceStart, "force", jobStartFlags.ForceStart,
		"Force stop jobs",
	)
	return cmd
}

func jobStart(client *apiclient.ApiClient, flags *JobStartFlags) error {
	filter := jobs.JobsFilterParams{
		Name:   flags.Name,
		AgeMin: flags.AgeMin,
		AgeMax: flags.AgeMax,
		Status: flags.Status,
	}

	var jobs jobs.Jobs
	err := jobs.GetFilteredJobs(client, &filter)
	if err != nil {
		return err
	}
	if len(jobs.Jobs) == 0 {
		return errors.New("no job matches your rules")
	}

	var action string
	if flags.Cron != "" {
		action = "schedule"
		fmt.Println("\nJobs to be scheduled :")
	} else {
		action = "start"
		fmt.Println("\nJobs to be started :")
	}
	jobs.PrintJobsTable()
	if !flags.ForceStart {
		err = askUserForYesOrNo(action)
		if err != nil {
			return err
		}
	}
	if action == "schedule" {
		return jobs.Schedule(client, flags.Cron)
	}
	return jobs.Start(client)
}
