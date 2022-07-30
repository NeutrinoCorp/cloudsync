# Neutrino CloudSync

![Go Build](https://github.com/NeutrinoCorp/cloudsync/actions/workflows/go.yml/badge.svg)
[![GoDoc](https://pkg.go.dev/badge/github.com/neutrinocorp/cloudsync)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/neutrinocorp/cloudsync)](https://goreportcard.com/report/github.com/neutrinocorp/cloudsync)
[![codebeat badge](https://codebeat.co/badges/2925d1e2-dbe0-4571-ba3c-0752db2b7e48)](https://codebeat.co/projects/github-com-neutrinocorp-cloudsync-master)
[![Coverage Status][cov-img]][cov]
[![Go Version][go-img]][go]

`Neutrino CloudSync` is an open-source tool used to upload entire file folders from any host to any cloud.

- [Neutrino CloudSync](#neutrino-cloudsync)
    - [Cloud Storage Drivers](#cloud-storage-drivers)
    - [How It Works](#how-it-works)
    - [Prerequisites](#prerequisites)
    - [How-To](#how-to)
        - [Provision your own infrastructure](#provision-your-own-infrastructure)
            - [Amazon Web Services](#amazon-web-services)
        - [Install](#install)
            - [Download binaries](#download-binaries)
            - [Running CLI application](#running-cli-application)
            - [Install globally](#install-globally)
        - [Update Configuration](#update-configuration)
        - [Upload Files (using compiled binary file)](#upload-files-using-compiled-binary-file)
        - [Upload Files (using source files)](#upload-files-using-source-files)

## Cloud Storage Drivers

Currently `CloudSync` offers integration with the following cloud storages:

- Amazon Simple Storage Service (S3)

And plans to add the following storages in a near future:

- Google Cloud Storage
- Microsoft Azure Blob Storage
- Google Drive

## How It Works

`CloudSync` comes with a CLI tool to execute operations while using cloud provider's Go SDK to interact with live
infrastructure.

## Prerequisites

- Go 1.18+
- Terraform
- AWS IAM user credentials configured with enough permissions to create/update:
    - S3 bucket
    - KMS key and alias
    - IAM user

**Optional**

- Make
- AWS CLI

_NOTE:_ `CloudSync` is able to use pre-existing cloud infrastructure if desired.

If that is the case, please consider provisioning a new IAM user with
enough roles/policies to interact with the blob storage.

## How-To

### Provision your own infrastructure

This repository contains fully-customizable `Terraform` code (IaaC) to provision required live infrastructure
in your own cloud account.

The code may be found [here](deployments/terraform).

To run this code and provision your own infrastructure, you MUST have installed `Terraform` CLI in your admin
host machine (not actual nodes which will interact with the cloud storage). Furthermore, an S3 bucket and a DynamoDB
table is REQUIRED to persist terraform states remotely (S3) and lock/unlock a remote mutex lock mechanism (DynamoDB),
enabling collaboration between multiple developers and hence development-purpose machines.
If this functionality is not desired, please remove the `main.tf`'s terraform block and leave it like this:

```terraform
terraform {
}
```

#### Amazon Web Services

The following steps are specific for the _Amazon Web Services (AWS)_ cloud provider:

- Go to deployments/terraform/workspaces/development.
- Add a `terraform.tfvars` file with the following variables (replace with actual cloud account data):

```text
aws_account = "0000"
aws_region = "us-east-N"
aws_access_key = "XXXX"
aws_secret_key = "XXXX"
```

*_The IAM user provided requires enough permissions to provision S3 buckets, IAM users and KMS keys with their alias._

- OPTIONAL: Modify variables from `variables.tf` file as desired to configure your infrastructure properties.
- Run the Terraform command `terraform plan` and verify a blob bucket and an encryption key will be created.
- Run the Terraform command `terraform apply` and write `yes` after reviewing all resources to be created.
- OPTIONAL: Go to the GUI cloud console (or use the cloud CLI) and verify all resources have been created
  with their proper configurations.

NOTE: At this time, the deployed infrastructure will get tagged and named using _development_ stage. This may be
removed through Terraform files, more specifically in the `main.tf` file from `development` workspace folder.

### Install

#### Download binaries

`Neutrino CloudSync` compiled binaries are available on the [releases page][releases] _(under Assets dropdown)_.

Download the binary file according your machine OS (Operating System, e.g. Windows, Mac/Darwin or Linux) and CPU
architecture (amd64, arm64).

#### Running CLI application

The binary file is a CLI program as a matter of fact, so it MUST be run within a terminal.

_Example:_

Linux/Darwin

```shell
user@machine:~ ./cloudsync -h
```

Windows (Powershell)

```shell
PS C:\Users\aruizeac> .\cloudsync.exe -h
```

#### Install globally

The binary file may be used as a standalone executable. Nevertheless, there is the option to install the executable so
it may run anywhere using a terminal.

In order to achieve this, move the binary file to user's homepath:

Linux/Darwin: /home/{USERNAME}/.cloudsync

Windows: C:\Users\\{USERNAME}\\.cloudsync

Finally, add the previous path to the _$PATH_ environment variable.

After compeleting all previous steps, the CLI application may be run like this:

Linux/Darwin

```shell
user@machine:~ cloudsync -h
```

Windows (Powershell)

```shell
PS C:\Users\aruizeac> cloudsync -h
```

_Notice the program does not require `.exe` nor `./` characters anymore._

### Update Configuration

`CloudSync` will create a new configuration file when running an operation _(e.g. upload command)_.

This file will be created under user's homepath inside a folder named _.cloudsync_. (/home/{USERNAME}/.cloudsync in **Linux/Darwin**, C:\Users\\{USERNAME}\\.cloudsync in **Windows**).

*DO NOT forget to enable view secret files/folders feature to see this folder.

| Field                            |    Type     | Description                                                                                                          |
|----------------------------------|:-----------:|:---------------------------------------------------------------------------------------------------------------------|
| cloud.region                     |   string    | Infrastructure region location _(e.g. us-east-1, us-west-2, eu-central-1)_                                           |
| cloud.bucket                     |   string    | Blob storage bucket name                                                                                             |
| cloud.access_key                 |   string    | Cloud account access key used to interact with infrastructure                                                        |
| cloud.secret_key                 |   string    | Cloud account access secret key used to interact with infrastructure                                                 |
| scanner.partition_id             |   string    | Identifier used to shard data within the blob storage _(auto-generated using ULID and might represent a machine ID)_ |
| scanner.read_hidden              |   boolean   | Enable scanning for hidden files                                                                                     |
| scanner.deep_traversing          |   boolean   | Enable scanning for child paths                                                                                      |
| scanner.ignored_keys             | string list | File or folder names to be ignored by scanner _(accepts wildcard patterns, e.g. *.go, *.java_)                       |
| scanner.log_errors               |   boolean   | Enable error logging                                                                                                 |

### Upload Files (using compiled binary file)

Run the `upload` command:

```shell
user@machine:~ cloudsync upload -d STORAGE_DRIVER -p DIRECTORY_TO_SYNC
```

For more information about the `upload` command, please run:

```shell
user@machine:~ cloudsync upload -h
```

### Upload Files (using source files)

Run the `cli` program using Go and execute `upload` command:

```shell
user@machine:~ go run ./cmd/cli/main.go upload -d STORAGE_DRIVER -p DIRECTORY_TO_SYNC
```

[actions]: https://github.com/neutrinocorp/cloudsync/workflows/Testing/badge.svg?branch=master

[godocs]: https://pkg.go.dev/github.com/neutrinocorp/cloudsync

[cov-img]: https://codecov.io/gh/NeutrinoCorp/cloudsync/branch/master/graph/badge.svg

[cov]: https://codecov.io/gh/NeutrinoCorp/cloudsync

[go-img]: https://img.shields.io/github/go-mod/go-version/NeutrinoCorp/cloudsync?style=square

[go]: https://github.com/NeutrinoCorp/cloudsync/blob/master/go.mod

[examples]: https://github.com/neutrinocorp/cloudsync/tree/master/examples

[releases]: https://github.com/NeutrinoCorp/cloudsync/releases
