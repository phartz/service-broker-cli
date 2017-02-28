# Service-Broker-CLI
Service Broker CLI is a commandline to to interact with a cloud foundry like service broker.
It was written to reduce to time to test changes on a service broker in the development phase.

## Download and install
Service broker cli is written in [`golang`](https://golang.org). It is designed to use only standard libraries.

Downlad the repo with `go get phartz.dedyn.io/gogs/phartz/service-broker-cli` then change the current path with cd `cd $GOPATH/src/phartz.dedyn.io/gogs/phartz/service-broker-cli`.

Get the depencies:
```
go get golang.org/x/crypto/ssh/terminal

go get github.com/fatih/color
```

Install with `go install service-broker-cli`and build with `go build -o sb`.

Now you can copy it to your bin folder `cp sb /usr/local/bin`

## Usage
```
$sb help

NAME:
   sb - A command line tool to interact with a service broker

USAGE:
   sb [global options] command [arguments...] [command options]

GETTING STARTED:
   help                                   Show help
   version                                Print the version
   status                                 Print the status information
   login                                  Log user in
   logout                                 Log user out
   target                                 Set or view the targeted org or space
   auth                                   User authentication

SERVICES:
   marketplace                            List available offerings in the marketplace
   services                               List all service instances in the target space
   service                                Show service instance info

   create-service                         Create a service instance
   update-service                         Update a service instance
   delete-service                         Delete a service instance

   create-service-key*                    Create key for a service instance
   service-keys*                          List keys for a service instance
   service-key*                           Show service key info
   delete-service-key*                    Delete a service key

```
## Log In

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
 