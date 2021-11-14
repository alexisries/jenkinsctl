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
	"jenkinsctl/pkg/apiclient"
	"jenkinsctl/pkg/apiclient/jobs"

	"github.com/spf13/cobra"
)

type JobListFlags struct {
	Name   string
	AgeMin int
	AgeMax int
	Status string
}

func newJobListFlags() *JobListFlags {
	return &JobListFlags{
		Name:   "",
		AgeMin: 0,
		AgeMax: 0,
		Status: "all",
	}
}

func NewJobListCmd(client *apiclient.ApiClient) *cobra.Command {
	jobListFlags := newJobListFlags()

	// cmd represents the job command
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list jobs",
		Long: `This command will list an jobs on Jenkins
For example:
	jenkinsctl job list
	jenkinsctl job list --state=running --minimum-age=30 --maximum-age=3600
	jenkinsctl job list --name=my-app`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return jobList(client, jobListFlags)
		},
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(
		&jobListFlags.Name, "name", jobListFlags.Name, "Specify the name of the job",
	)
	cmd.Flags().StringVar(
		&jobListFlags.Status, "status", jobListFlags.Status,
		"Filter Job from status (possible values: all, running, success, aborted, failure)",
	)
	cmd.Flags().IntVar(
		&jobListFlags.AgeMin, "minimum-age", jobListFlags.AgeMin,
		"Filter Jobs from last build minimum age (in minutes)",
	)
	cmd.Flags().IntVar(
		&jobListFlags.AgeMin, "maximum-age", jobListFlags.AgeMin,
		"Filter Jobs from last build maximum age (in minutes)",
	)
	return cmd
}

func jobList(client *apiclient.ApiClient, flags *JobListFlags) error {
	filter := jobs.JobsFilterParams{
		Name:   flags.Name,
		AgeMin: flags.AgeMin,
		AgeMax: flags.AgeMax,
		Status: flags.Status,
	}
	err := checkStatusValidValue(flags.Status)
	if err != nil {
		return err
	}
	var jobs jobs.Jobs
	err = jobs.GetFilteredJobs(client, &filter)
	if err != nil {
		return err
	}
	if len(jobs.Jobs) == 0 {
		return errors.New("no job matches your rules")
	}
	jobs.PrintJobsTable()
	return nil
}
