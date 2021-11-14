# Jenkins CLI client

Jenkinsctl is a command line client to manage a Jenkins server.

For the moment, only job management is implemented, so you can:
- list the jobs according to certain criteria such as the age and status of the last build or the name of the job.
- stop jobs with the same criteria as the listing
- start jobs immediately or by scheduling it with Jenkins schedule syntax

## Build with docker

To build the client from source without having to install the `golang` development environment, you can create a docker image. 

Start by building the docker image with `docker build` 

```shell
$ docker build -t jenkinsctl .
```

## Use docker image

To load the credentials, two solutions are available to you:
- Use environment variables
- Mount a volume in the container 

### Configuration in `environment variables`

To start the program with environment variables, run this command: 

```shell
$ docker run --rm -it --name jenkinsctl \
        -e JENKIN_ADDR=http://jenkins.local \
        -e JENKINS_USER=jenkins_user \
        -e JENKINS_TOKEN=jenkins_token \
        jenkinsctl --help
```

### Configuration in `.jenkinsctl.yaml`

Instead of using environment variables, you can use a `.jenkinsctl.yaml` configuration file.

```yaml
jenkins:
  addr: http://jenkins.local
  user: jenkins_user
  token: jenkins_token
```

For more details on the contents of the configuration file, go to the "Configuration" section. 

Place your `.jenkinsctl.yaml` file in the current directory and run the docker image with a mount of the current volume in the `/app` directory of the container :

```shell
$ docker run --rm -it --name jenkinsclient -v $(pwd):/app jenkinsctl --help
```

## Configuration

Below are the configuration options :

| Name                    | Description                                       | Default            |
| ----------------------- | ------------------------------------------------- | ------------------ |
| `jenkins.addr`          | Address of the Jenkins server                     | `http://localhost` |
| `jenkins.user`          | Jenkins account username                          | `""`               |
| `jenkins.token`         | API token of the Jenkins account                  | `""`               |
| `jenkins.max_concurent` | Maximum number of concurent http requests         | `3`                |

By default the program will read the configuration in the `.jenkinsctl.yaml` file at these paths: 

1. In the current working directory
2. In the user home directory

> **Tip**: You can also specify the config file with the `--config` flag

### Configuration via environment variables

You can use environment variables to supplement or replace the configuration file.

To do this, put the name of the configuration otions in uppercase and replace the dots with underscores, for example:

```shell
$ JENKINS_MAX_CONCURENT=10 ./jenkinsctl
```

## Commands

We will see here the differents commands of jenkinsctl

### List jobs

To list the jobs on the Jenkins server, you can use the `jenkinsctl job list` command 

#### Command flags (optional)

| Name            | Description                                                                                         | Default |
| --------------- | ----------------------------------------------------------------------------------------------------| ------- |
| `--name`        | Specify the name of the job you want to get                                                         | `""`    |
| `--status`      | Filter jobs from status of the last build (possible values: all, running, success, aborted, failure)| `all`   |
| `--minimum-age` | Filter jobs from last build minimum age (in minutes)                                                | `""`    |
| `--maximum-age` | Filter jobs from last build maximum age (in minutes)                                                | `""`    |

### Start jobs

To start the jobs on the Jenkins server, you can use the `jenkinsctl job start` command 


#### Command flags (optional)

| Name            | Description                                                                                         | Default |
| --------------- | ----------------------------------------------------------------------------------------------------| ------- |
| `--name`        | Specify the name of the job you want to start                                                       | `""`    |
| `--status`      | Filter jobs from status of the last build (possible values: all, running, success, aborted, failure)| `all`   |
| `--minimum-age` | Filter jobs from last build minimum age (in minutes)                                                | `""`    |
| `--maximum-age` | Filter jobs from last build maximum age (in minutes)                                                | `""`    |
| `--schedule`    | Specify the schedule in Jenkins time trigger syntax                                                 | `""`    |
| `--force`       | Do not ask for confirmation before starting                                                         | `false` |

### Stop jobs

To stop the jobs on the Jenkins server, you can use the `jenkinsctl job stop` command 


#### Command flags (optional)

| Name            | Description                                                                                         | Default |
| --------------- | ----------------------------------------------------------------------------------------------------| ------- |
| `--name`        | Specify the name of the job you want to stop                                                        | `""`    |
| `--status`      | Filter jobs from status of the last build (possible values: all, running, success, aborted, failure)| `all`   |
| `--minimum-age` | Filter jobs from last build minimum age (in minutes)                                                | `""`    |
| `--maximum-age` | Filter jobs from last build maximum age (in minutes)                                                | `""`    |
| `--force`       | Do not ask for confirmation before stopping                                                         | `false` |


## Examples

We will see here the different possibilities offered by this program 

### lists all running jobs on a Jenkins server

To list all running jobs on the Jenkins server, start the `job list` command with `--status running` flag :

```shell
$ jenkinsctl job list --status running
+--------------+---------+---------------------+
|     NAME     | STATUS  |     BUILD DATE      |
+--------------+---------+---------------------+
| test-java    | running | 2021-11-14 12:06:45 |
| test-java-10 | running | 2021-11-14 12:06:46 |
| test-java-5  | running | 2021-11-14 12:06:49 |
| test-java-8  | running | 2021-11-14 12:06:51 |
+--------------+---------+---------------------+
```

### stops all jobs that have been running for more than 1 hour

To stop all jobs that are being executed for over an hour, launch the `job stop` command with `--minimum-age 60` flag :

```shell
$ jenkinsctl job stop  --minimum-age 60

Jobs to be stopped :
+-------------+---------+---------------------+
|    NAME     | STATUS  |     BUILD DATE      |
+-------------+---------+---------------------+
| test-java-5 | running | 2021-11-14 12:06:49 |
| test-java-8 | running | 2021-11-14 12:06:51 |
+-------------+---------+---------------------+

Do you want to stop these jobs ? (yes or no): yes

Stopping jobs...
job test-java-5 is now in stopped state
job test-java-8 is now in stopped state
```

### starts a job at a given hour

To start a job at 8am every day, lauch `job start` command with `--name yourjob` and `--schedule H 8 * * *` flags:

```shell
$ jenkinsctl job start --name "test-java-5"  --schedule "H 8 * * *"

Jobs to be scheduled :
+-------------+---------+---------------------+
|    NAME     | STATUS  |     BUILD DATE      |
+-------------+---------+---------------------+
| test-java-5 | aborted | 2021-11-14 12:06:49 |
+-------------+---------+---------------------+

Do you want to schedule these jobs ? (yes or no): yes

Scheduling jobs...
job test-java-5 is now scheduled (H 8 * * *)
```