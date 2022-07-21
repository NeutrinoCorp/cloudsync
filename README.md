# Neutrino CloudSync

![Go Build](https://github.com/NeutrinoCorp/cloudsync/actions/workflows/go.yml/badge.svg)
[![GoDoc](https://pkg.go.dev/badge/github.com/neutrinocorp/cloudsync)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/neutrinocorp/cloudsync)](https://goreportcard.com/report/github.com/neutrinocorp/cloudsync)
[![codebeat badge](https://codebeat.co/badges/2925d1e2-dbe0-4571-ba3c-0752db2b7e48)](https://codebeat.co/projects/github-com-neutrinocorp-cloudsync-master)
[![Coverage Status][cov-img]][cov]
[![Go Version][go-img]][go]

`Neutrino CloudSync` is an open-source tool used to upload entire file folders from any host to any cloud.

- [Neutrino CloudSync](#neutrino-cloudsync)
    - [How-To](#how-to)
        - [Provision your own infrastructure](#provision-your-own-infrastructure)
            - [Amazon Web Services](#amazon-web-services)
        - [Upload Files](#upload-files)

## How-To

### Provision your own infrastructure

This repository contains fully-customizable `Terraform` code (IaC) to deploy and/or provision live infrastructure
in your own cloud account.

The code should be found [here](deployments/terraform).

To run this code and provision your own infrastructure, you MUST have installed `Terraform` CLI in your admin
host machine (not actual nodes which will interact with the platform's). Furthermore, an S3 bucket and a DynamoDB
table is REQUIRED to persist terraform states remotely (S3) and lock/unlock a remote mutex lock mechanism (DynamoDB),
enabling collaboration between multiple developers and hence development-purpose machines.
If this functionality is not desired, please remove the `main.tf`'s terraform block and leave it like this:

```terraform
terraform {
}
```

#### Amazon Web Services

The following steps are specific for the _Amazon Web Services (AWS)_ platform:

- Go to deployments/terraform/workspaces/development.
- Add a `terraform.tfvars` file with the following variables (replace with actual cloud account data):

```text
aws_account = "0000"
aws_region = "us-east-N"
aws_access_key = "XXXX"
aws_secret_key = "XXXX"
```

- OPTIONAL: Modify variables from `variables.tf` file as desired to configure your infrastructure properties.
- Run the Terraform command `terraform plan` and verify a blob bucket and an encryption key will be created.
- Run the Terraform command `terraform apply` and write `yes` after verifying a blob bucket and an encryption key are
  the only resources to be created.
- OPTIONAL: Go to the GUI cloud console (or use the cloud CLI) and verify all resources have been created
  with their proper configurations.

NOTE: At this time, the deployed infrastructure will get tagged and named using _development_ stage. This may be
removed through Terraform files, more specifically in the `main.tf` file from `development` workspace folder.

### Upload Files

Just start an `uploader` instance using the command:

```shell
user@machine:~ make run directory=DIRECTORY_TO_SYNC
```

or

```shell
user@machine:~ go run ./cmd/uploader/main.go -d DIRECTORY_TO_SYNC
```

The uploader program has other flags not specified on the examples. To read more about them please run the command:

```shell
user@machine:~ make help
```

or

```shell
user@machine:~ go run ./cmd/uploader/main.go -h
```

[actions]: https://github.com/neutrinocorp/cloudsync/workflows/Testing/badge.svg?branch=master
[godocs]: https://pkg.go.dev/github.com/neutrinocorp/cloudsync
[cov-img]: https://codecov.io/gh/NeutrinoCorp/cloudsync/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/NeutrinoCorp/cloudsync
[go-img]: https://img.shields.io/github/go-mod/go-version/NeutrinoCorp/cloudsync?style=square
[go]: https://github.com/NeutrinoCorp/cloudsync/blob/master/go.mod
[examples]: https://github.com/neutrinocorp/cloudsync/tree/master/examples
