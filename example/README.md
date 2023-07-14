# Example

This deploys a single openid client.

## Before you run.

Since this provider is not yet registered with hashicorp, you need to override the provider isntallation
to allow it to find the compiled provider.

Make something like this be your .terraformrc or terraform.rc file:

```code
#
# ~/.terraformrc configuration file.
#
provider_installation {
  filesystem_mirror {
    path = "/Users/brad/.terraform.d/plugins"
    include = ["local.providers/*/*"]
  }
  direct {
      exclude = ["local.providers/*/*"]
  }
}
```

## Initialize

```code
terraform init
```

### Build the Plan

```code
terraform plan
```

## Apply the Plan

```code
terraform apply
```

Now goof around with the configuration files: add another client, modify a client, remove a line and
re-run the plan and apply and you'll see the state mimic the changes.

