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
	"fmt"
	"jenkinsctl/pkg/apiclient"
)

func (jobs *Jobs) Start(clt *apiclient.ApiClient) error {
	fmt.Println("Starting jobs...")
	for _, job := range jobs.Jobs {
		if job.IsRunning {
			fmt.Printf("job %s is already in started state\n", job.Name)
			continue
		}
		buildId, err := job.JenkinsJob.InvokeSimple(clt.Ctx, map[string]string{})
		if err != nil {
			return err
		}
		if buildId > 0 {
			fmt.Printf("job %s is now in started state, build id: %d\n", job.Name, buildId)
		}
	}
	return nil
}
