# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

locals {
  smilecdr_version      = "2021.05.R01"
  callback_script       = file("${path.module}/js/local_inbound_security/authentication_callback.js")
  auth_script_1         = file("${path.module}/js/idp_patients/on_auth_success.js")
  auth_script_2         = file("${path.module}/js/idp_patients/on_auth_success.js")
  user_mapping_script = file("${path.module}/js/idp_patients/user_mapping.js")
  example_jwks = file("${path.module}/jwks/example.jwks.json")
  post_authorize_callback = file("${path.module}/js/smart_outbound_security/post_authorize_callback.js")
  codap_authorization_callback = file("${path.module}/js/smart_outbound_security/codap_authorization.js")
}
