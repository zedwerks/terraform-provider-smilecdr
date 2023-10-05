# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

locals {
  smilecdr_version             = "2021.05.R01"
  inbound_callback_script      = file("${path.module}/js/local_inbound_security/authentication_callback.js")
  idp_auth_script              = file("${path.module}/js/idp/on_auth_success.js")
  user_mapping_script          = file("${path.module}/js/idp/user_mapping.js")
  example_jwks                 = file("${path.module}/jwks/example.jwks.json")
  post_authorize_callback      = file("${path.module}/js/smart_outbound_security/post_authorize_callback.js")
  codap_authorization_callback = file("${path.module}/js/smart_outbound_security/codap_authorization.js")
}
