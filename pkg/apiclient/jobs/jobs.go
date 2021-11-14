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

package jobs

import (
	"os"
	"strings"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/olekukonko/tablewriter"
)

const (
	JOB_STATUS_ALL     = "all"
	JOB_STATUS_RUNNING = "running"
	JOB_STATUS_SUCCESS = "success"
	JOB_STATUS_FAILED  = "failure"
	JOB_STATUS_ABORTED = "aborted"
	JOB_STATUS_NOBUILD = "no_build"
)

type Job struct {
	Id                    int
	Name                  string
	LastBuildDuration     float64
	LastBuildCreationDate time.Time
	IsRunning             bool
	Success               bool
	Result                string
	JenkinsJob            *gojenkins.Job
	JenkinsLastBuild      *gojenkins.Build
}

type Jobs struct {
	Jobs []Job
}

type JobsFilterParams struct {
	Name   string
	AgeMin int
	AgeMax int
	Status string
}

func (job *Job) checkJobStatusMatch(status string) bool {
	switch {
	case status == JOB_STATUS_ALL:
		return true
	case status == JOB_STATUS_RUNNING && job.IsRunning:
		return true
	case status == JOB_STATUS_SUCCESS && !job.IsRunning && job.Success:
		return true
	case status == JOB_STATUS_FAILED && !job.IsRunning && !job.Success:
		return true
	case status == JOB_STATUS_ABORTED && job.Result == gojenkins.STATUS_ABORTED:
		return true
	default:
		return false
	}
}

func (job *Job) checkJobMinimumAgeMatch(ageMin int) bool {
	if ageMin == 0 {
		return true
	}
	diff := time.Since(job.LastBuildCreationDate).Minutes()
	return diff >= float64(ageMin)
}

func (job *Job) checkJobMaximumAgeMatch(ageMax int) bool {
	if ageMax == 0 {
		return true
	}
	diff := time.Since(job.LastBuildCreationDate).Minutes()
	return diff < float64(ageMax)
}

func (jobs *Jobs) PrintJobsTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Build date"})
	for _, job := range jobs.Jobs {
		var jobStatus string
		if job.IsRunning {
			jobStatus = "running"
		} else {
			jobStatus = strings.ToLower(job.Result)
		}

		var jobBuildDateStr string
		if job.Result == JOB_STATUS_NOBUILD {
			jobBuildDateStr = ""
		} else {
			jobBuildDateStr = job.LastBuildCreationDate.Format("2006-01-02 15:04:05")
		}

		j := []string{
			job.Name,
			jobStatus,
			jobBuildDateStr,
		}
		table.Append(j)
	}
	table.Render()
}
