# Terraform Provider for Smile CDR

![Badge-License]

An initial cut at creating a Terraform provider for Smile CDR. Focus early days on security module configurations:

- OpenID Identity Providers
- OIDC Clients
- Authorization and Permissions.

Then later on, add more comprehensive configurations for Smile CDR.

## Build

Run the following command to build the provider

```shell
go build -o terraform-provider-smilecdr
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the example configuration.

```shell
terraform init && terraform apply
```

[Badge-License]: https://img.shields.io/badge/license-apache%202.0-60C060.svg


## Running Acceptance Tests

```code
make testacc
```
