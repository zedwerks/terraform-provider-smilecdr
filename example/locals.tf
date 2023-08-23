# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

locals {
  smilecdr_version      = "2021.05.R01"
  callback_script       = file("${path.module}/js/inbound_security_callback.js")
  auth_script_1         = file("${path.module}/js/federationAuthScript.js")
  auth_script_2         = file("${path.module}/js/federationAuthScript.js")
  user_mapping_script_1 = file("${path.module}/js/federationUserMappingScript.js")
  user_mapping_script_2 = file("${path.module}/js/federationUserMappingScript.js")
  example_jwks = file("${path.module}/jwks/example.jwks.json")
}
