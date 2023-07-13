# Example

This deploys a single openid client.

## Before you run.

Since this provider is not yet registered with hashicorp, you need to override the provider isntallation
to allow it to find the compiled provider.

Make something like this be your .terraformrc or terraform.rc file:

```code
#
# .terraformrc configuration file.
#
provider_installation {
  dev_overrides {
    "registry.terraform.io/zedwerks/smilecdr" = "/Users/brad/.terraform.d/plugins/zedwerks/smilecdr"
  }
}
```
