# Terraform Provider for Smile CDR

![Badge-License]

A a Terraform provider for Smile CDR. Focus on security module configurations:

- OpenID Identity Providers
- OIDC Clients
- Authorization and Permissions.
- User Accounts
- SMART Outbound Security Modules
- SMART Inbound Security Module

## Requirements

- Terraform >=10.0
- Go >= 1.19

## Installing the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the GO install command:

```shell
make install
```

## Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the GO install command:

```shell
make install
```

Run the following command to build the provider

```shell
make build
```

## Using the Provider

Documentation for this Terraform provider can be found [here](./docs).


1. Have Smile CDR installed, possibly locally. See the smilecdr-product/docker-compose.yml, if you have a docker image for a release of Smile CDR and you wish to run locally.
2. Install the provider
3. Have Smile CDR installed, possibly locally. See the smilecdr-product/docker-compose.yml, if you have a docker image for a release of Smile CDR and you wish to run locally.
4. Install the provider

```shell
make install
```

3. Change directory to the example folder
4. Set three environment variables:
   - SMILECDR_USERNAME
   - SMILECDR_PASSWORD
   - SMILECDR_BASE_URL
5. 3. Change directory to the example folder
4. Set three environment variables:
   - SMILECDR_USERNAME
   - SMILECDR_PASSWORD
   - SMILECDR_BASE_URL
5. Then, run the following command to initialize the workspace and apply the example configuration.

```shell
terraform init && terraform apply
```

## Running the Acceptance Tests

To run acceptance tests you will need the following environment variables set so that the acceptance tests can connect to a dev/test instance of Smile CDR:

- `SMILECDR_USERNAME`, which is an admin user, that has full access to Admin APIs.

- `SMILECDR_PASSWORD`, the password for the admin user.

- `SMILECDR_BASE_URL`, which is the base URL of the Smile CDR Admin APIs

*Note:* Acceptance tests create real resources, and often cost money to run. That said, my recommendation is
to run your acceptance tests on a local desktop Docker image. See the smilecdr-product folder for a starting-point docker compose YAML file.


```shell
make testacc
```

## Working on this Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

## Using Goreleaser

This project uses [goreleaser](goreleaser.com) to manage publishing signed release binaries to Git. From there they are published to the Terraform Registry (terraform.io)

See ```goreleaser.com/quick-start```

## Creating a New Release with GitActions setup (as per this repo)

1. Merge from working branch to main.
2. Tag a new release in form "vx.y.z" e.g. v1.0.1. OPtionally add a message with -m

```shell
git tag -a v1.0.5 -m "Release 1.0.5" 
git push origin v1.0.5
```

3. This triggers the GitAction to run a deployment.
4. Monitor the GitAction progress at Github.com

[Badge-License]: https://img.shields.io/badge/license-apache%202.0-60C060.svg
