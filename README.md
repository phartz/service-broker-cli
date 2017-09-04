# Service-Broker-CLI
Service Broker CLI is a commandline tool to interact with a cloud foundry like service broker.
It was written to reduce to time to test changes on a service broker in the development phase.

## Table of Content
This documentation contains the following topics.

* [Download, Build and Install](##Download-Build-and-Install)
* [Usage](##Usage)
* [Restrictions](##Restrictions)
* [Log In](##Log-In)
* [Logging](##Logging)

## Download, Build and Install
Service broker cli is written in [`golang`](https://golang.org). It is designed to use only standard libraries.

It provides also a Makefile, so it's quite easy to build and install.

Downlad the repo with `go get github.com/phartz/service-broker-cli` then change the current path to the sources `cd $GOPATH/src/github.com/phartz/service-broker-cli`.

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

    target                t           Sets or gets the target
    login                 l           Login to the target
    logout                lo          Logout from the target
    auth                              Authenticate to the target
    version               -v          Print the version

    marketplace           m           List available offerings in the marketplace
    services              s           List all service instances in the target space
    service                           Show service instance info

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
* The service broker doesn't store any user crendetials so it is not possible to get information about a service key.
In contrast to to CloudFounndry CLI the `sb create-service-key` returns the credentials.
* The service broker also doesn't store any information about the real service names, the broker only works with UUIDs.


## Target and log In

Use `sb target` to target a Service Broker.

```bash
$ sb target http://localhost:3000
```

To log in you can use either `sb login` as an interactive operation or `sb auth` for scripting.

The credentials will be stored in json format in the file `.sb`. Those file can be located either in the current working directory or in its parent directory. If the file was not found, sb is looking also in the users root path.

```
--> /some/where/on/your/disk
--> /some/where/on/your
--> /some/where/on
--> /some/where
--> /some
--> /
--> ~
```

## Logging

Service-Broker-CLI does not provide a logging in the classical sense. But you can trace the Service-Broker API requests and answers.

```
$ SB_TRACE=ON && sb services
```
