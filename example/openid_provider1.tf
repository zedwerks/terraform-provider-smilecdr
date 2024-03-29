# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_openid_identity_provider" "example_idp1" {
  depends_on = [ smilecdr_smart_outbound_security.smart_auth_federated ]
  node_id                             = "Master"
  module_id                           = "smart_auth_federated"
  issuer                              = "https://example1.com/auth/issuer"
  name                                = "example_idp1"
  federation_request_scopes           = "openid profile email"
  federation_authorization_url        = "http://localhost:8800/auth"
  federation_token_url                = "http://localhost:8800/token"
  federation_user_info_url            = "http://localhost:8800/userinfo"
  federation_jwk_set_url              = "http://localhost:8800/auth/jwks"
  federation_auth_script_text         = local.idp_auth_script
  token_introspection_client_id       = "smile"
  token_introspection_client_secret   = "client_secret_goes_here"
}
