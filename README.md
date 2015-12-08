check_kube
============

A Nagios check to verify the status of a Kubernetes cluster resources

Installing
----------
If you would rather not build the binaries yourself, you can install compiled,
statically-linked [binaries](https://github.com/jfchevrette/check_kube/releases)

Building
--------
```
go get github.com/jfchevrette/check_kube

# Binary will be at $GOPATH/bin/check_kube
```

Usage
-----
```
NAME:
   check_kube_nodes - Nagios check to verify Kubernetes resources status

USAGE:
   check_kube [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   node, n	check node status
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-endpoint 	Kubernetes API Endpoint
   --username 		Kubernetes API Username
   --password 		Kubernetes API Password
   --skip-tls-verify	Skip TLS certificate verification
   --help, -h		show help
   --version, -v	print the version
```

`--api-endpoint`: Here you specify your Kubernetes API endpoint URL.

`--username`: Kubernetes API Username for basic authentication

`--password`: Kubernetes API Password for basic authentication

`--skip-tls-verify`:Skip TLS certificate verification when you have a self-signed certificate.

Contributions
-------------
Contributions are more than welcome. Bug reports with specific reproduction
steps are great. If you have a code contribution you'd like to make, open a
pull request with suggested code.
