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
	"errors"
	"jenkinsctl/pkg/apiclient"
	"sort"
	"sync"
	"time"

	"github.com/bndr/gojenkins"
)

type JenkinsJobNum struct {
	Id  int
	Job *gojenkins.Job
}

func (job *Job) parseJenkinsJobWithBuilds(
	clt *apiclient.ApiClient, jenkinsJob JenkinsJobNum,
) error {
	job.Id = jenkinsJob.Id
	job.Name = jenkinsJob.Job.GetName()
	job.JenkinsJob = jenkinsJob.Job

	last_build, err := jenkinsJob.Job.GetLastBuild(clt.Ctx)
	if err != nil {
		if err.Error() == "404" {
			job.LastBuildDuration = 0
			job.LastBuildCreationDate = time.Time{}
			job.IsRunning = false
			job.Result = JOB_STATUS_NOBUILD
			job.JenkinsLastBuild = &gojenkins.Build{}
			return nil
		} else {
			return err
		}
	}
	job.LastBuildDuration = last_build.GetDuration()
	job.LastBuildCreationDate = last_build.GetTimestamp()
	job.IsRunning = last_build.Raw.Building
	job.Result = last_build.GetResult()
	job.JenkinsLastBuild = last_build
	return nil
}

func (job *Job) getJobByName(clt *apiclient.ApiClient, name string) error {
	jenkinsJob, err := clt.Jenkins.GetJob(clt.Ctx, name)
	if err != nil {
		if err.Error() == "404" {
			return errors.New("job not found")
		}
		return err
	}

	err = job.parseJenkinsJobWithBuilds(clt, JenkinsJobNum{Id: 1, Job: jenkinsJob})
	if err != nil {
		return err
	}
	return nil
}

func (jobs *Jobs) getAllJobs(clt *apiclient.ApiClient) error {
	jenkinsJobs, err := clt.Jenkins.GetAllJobs(clt.Ctx)
	if err != nil {
		return err
	}
	jenkinsJobsNum := []JenkinsJobNum{}
	for i, jenkinsJob := range jenkinsJobs {
		jenkinsJobsNum = append(jenkinsJobsNum, JenkinsJobNum{Id: i, Job: jenkinsJob})
	}

	doGetJobsBuilds := make(chan JenkinsJobNum, len(jenkinsJobsNum))
	for _, jenkinsJob := range jenkinsJobsNum {
		doGetJobsBuilds <- jenkinsJob
	}
	close(doGetJobsBuilds)

	jobsWithBuilds := make(chan Job, len(jenkinsJobsNum))
	var wg sync.WaitGroup
	for i := 0; i < clt.ClientConfig.MaxConcurentRequests; i++ {
		wg.Add(1)
		go func() {
			for jenkinsJobNum := range doGetJobsBuilds {
				job := Job{}
				err := job.parseJenkinsJobWithBuilds(clt, jenkinsJobNum)
				if err != nil {
					panic(err)
				}
				jobsWithBuilds <- job
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(jobsWithBuilds)

	for job := range jobsWithBuilds {
		jobs.Jobs = append(jobs.Jobs, job)
	}
	sort.Slice(jobs.Jobs, func(i, j int) bool {
		return jobs.Jobs[i].Id < jobs.Jobs[j].Id
	})
	return nil
}

func (jobs *Jobs) GetFilteredJobs(
	clt *apiclient.ApiClient, filter *JobsFilterParams,
) error {
	if filter.Name != "" {
		job := Job{}
		err := job.getJobByName(clt, filter.Name)
		if err != nil {
			return err
		}
		jobs.Jobs = []Job{job}
		return nil
	}
	jobsInput := Jobs{}
	err := jobsInput.getAllJobs(clt)
	if err != nil {
		return err
	}

	for _, job := range jobsInput.Jobs {
		if !job.checkJobStatusMatch(filter.Status) ||
			!job.checkJobMinimumAgeMatch(filter.AgeMin) ||
			!job.checkJobMaximumAgeMatch(filter.AgeMax) {
			continue
		}
		jobs.Jobs = append(jobs.Jobs, job)
	}
	return nil
}
