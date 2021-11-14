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

	"github.com/bndr/gojenkins"
	"github.com/spf13/viper"
)

type ApiClient struct {
	Jenkins      *gojenkins.Jenkins
	Ctx          context.Context
	ClientConfig ApiClientConfig
}

type ApiClientConfig struct {
	MaxConcurentRequests int
}

func (clt *ApiClient) Initialize() {
	viper.SetDefault("jenkins.addr", "localhost")
	addr := viper.GetString("jenkins.addr")

	username := viper.GetString("jenkins.user")
	password := viper.GetString("jenkins.token")

	viper.SetDefault("jenkins.max_concurent", 3)
	clt.ClientConfig = ApiClientConfig{
		MaxConcurentRequests: viper.GetInt("jenkins.max_concurent"),
	}
	clt.Ctx = context.Background()
	clt.Jenkins = gojenkins.CreateJenkins(nil, addr, username, password)
	_, err := clt.Jenkins.Init(clt.Ctx)
	if err != nil {
		panic(err)
	}
}
