# terraform-provider-smilecdr Change Log

## v1.0.5 (Dec 21, 2023)

- Add import support for resources, except for new ```smilecdr_user``` resource.
- Suppressed plan changes when managing user resource passwords.
- Added guidance to ignore changes for ```smilecdr_openid_client``` ```client_secrets``` due to complex masking and other changes.
- Updated acceptance tests for above.
- Updated docs to include import and new resource ```smilecdr_user```

Since many modules are pre-built in SmileCDR,such as the ```local_security``` and ```smart_auth``` modules, you can now import those modules and manage their configurations with
Terraform.

This release adds user management but this resource has only been tested against an inbound security module. It is useful for bootstrapping test accounts in a non-production environment.

## v1.0.4 (Nov 15, 2023)

- attempt to fix goreleaser build for binaries to resolve not found errors.

## v1.0.2 (Oct 26, 2023)

- updated introductory documentation

## v1.0.1 (Oct 26, 2023)

- updated gorelease build to factor multi-archs per os.

## v1.0.0 (Sep 27, 2023)

FEATURES:

- configure OIDC clients
- configure OIDC Federated Authorization Servers
- configure Outbound Security Module
- configure generic module (experimental)
