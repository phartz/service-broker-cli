# Service Broker CLI
Service Broker CLI is a command-line tool to interact with a Cloud Foundry like Service Broker.
It was written to reduce to time to test changes on a Service Broker in the development phase.

## Table of Contents
This documentation contains the following topics.

* [Download, Build and Install](##Download-Build-and-Install)
* [Usage](##Usage)
* [Restrictions](##Restrictions)
* [Log In](##Log-In)
* [Logging](##Logging)
* [CI Pipeline](##ci-pipeline)

## Download, Build and Install
Service Broker CLI is written in [`Go`](https://golang.org). It is designed to use only standard libraries.

It provides also a Makefile, so it's quite easy to build and install.

Downlad the repo with `go get github.com/phartz/service-broker-cli`, then change the current path to the sources `cd $GOPATH/src/github.com/phartz/service-broker-cli`.

Build
```
make
```

Install
```
make install
```

Either add the `$GOPATH/bin` `export PATH=$PATH:$GOPATH/bin` to your `$PATH` or copy the cli to your bin folder `cp sb /usr/local/bin`

## Usage
```
$sb help

NAME:
   sb - A command line tool to interact with a service broker

USAGE:
   sb [global options] command [arguments...] [command options]

COMMANDS

    COMMAND               SHORTCUT    DESCRIPTION
    help                  h           Show help

    api                               Set or view target api url
    login                 l           Login to the target
    logout                lo          Logout from the target
    auth                              Authenticate to the target
    version               -v          Print the version

    marketplace           m           List available offerings in the marketplace
    services              s           List all service instances in the target space
    service                           Show service instance info
    find-service          fs          Lookup a service with a given deployment name

    create-service        cs          Create a service instance
    update-service                    Update a service instance
    delete-service        ds          Delete a service instance

    create-service-key    csk         Create key for a service instance
    service-keys          sk          List keys for a service instance
    delete-service-key    dsk         Delete a service key

```

To get specific help use `sb help <command>`


## Restrictions

There are some restrictions:
* The Service Broker doesn't store any user credentials so it is not possible to get information about a service key.
In contrast to to Cloud Foundry CLI the `sb create-service-key` returns the credentials.
* The Service Broker also doesn't store any information about the real service names. The broker only works with UUIDs.


## Target and log In

Use `sb target` to target a Service Broker.

```bash
$ sb target http://localhost:3000
```

To log in you can use either `sb login` as an interactive operation or `sb auth` for scripting.

The credentials will be stored in JSON format in the file `.sb`. Those files can be located either in the current working directory or in its parent directory. If the file was not found, sb is looking also in the users root path.

```
--> /some/where/on/your/disk
--> /some/where/on/your
--> /some/where/on
--> /some/where
--> /some
--> /
--> ~
```

### Alternative Way to Use the CLI

The Service-Broker-CLI is able to read the credentials and the host from environment variables. No login is then needed.

```
$ SB_HOST="http://redis-service-broker.service.dc1.consul:3000/" SB_USERNAME="admin" SB_PASSWORD="mytopsecretpwd" sb ......
```

| Env Var | Description |
|---|---|
| SB_HOST | The host name, it can be the IP, the FQHN, with or without the scheme and port |
| SB_USERNAME | The username |
| SB_PASSWORD | The password |


## Logging

Service-Broker-CLI does not provide a logging in the classical sense. But you can trace the Service Broker API requests and answers.

```
$ SB_TRACE=ON && sb services
```

## CI Pipeline

A [concourse](concourse.ci) based CI pipeline can be found [here](https://phartz.dedyn.io/teams/main/pipelines/service-broker-cli)
