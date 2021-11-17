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

package apiclient

import (
	"context"
	"errors"

	"github.com/bndr/gojenkins"
	"github.com/spf13/viper"
)

type ApiClient struct {
	Jenkins      *gojenkins.Jenkins
	Ctx          context.Context
	ClientConfig *ApiClientConfig
}

type ApiClientConfig struct {
	address              string
	username             string
	token                string
	MaxConcurentRequests int
}

func (clt *ApiClient) getConfig() error {

	viper.SetDefault("jenkins.max_concurent", 3)

	clt.ClientConfig = &ApiClientConfig{
		address:              viper.GetString("jenkins.addr"),
		username:             viper.GetString("jenkins.user"),
		token:                viper.GetString("jenkins.token"),
		MaxConcurentRequests: viper.GetInt("jenkins.max_concurent"),
	}
	missingConfig := false
	var errorMessage = "\n"
	if clt.ClientConfig.address == "" {
		missingConfig = true
		errorMessage = errorMessage + "jenkins server address not defined\n"
	}
	if clt.ClientConfig.username == "" {
		missingConfig = true
		errorMessage = errorMessage + "jenkins server username not defined\n"
	}
	if clt.ClientConfig.token == "" {
		missingConfig = true
		errorMessage = errorMessage + "jenkins server token not defined\n"
	}
	if missingConfig {
		return errors.New(errorMessage)
	}
	return nil
}

func (clt *ApiClient) Initialize() {
	if err := clt.getConfig(); err != nil {
		panic(err)
	}
	clt.Ctx = context.Background()
	clt.Jenkins = gojenkins.CreateJenkins(
		nil, clt.ClientConfig.address, clt.ClientConfig.username, clt.ClientConfig.token,
	)
	if _, err := clt.Jenkins.Init(clt.Ctx); err != nil {
		panic(err)
	}
}
