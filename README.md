# Neutrino CloudSync

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
