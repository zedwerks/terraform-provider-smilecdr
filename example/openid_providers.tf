# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_openid_identity_provider" "example_idp1" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  issuer                              = "https://example1.com/auth/issuer"
  name                                = "example_idp1"
  federation_request_scopes           = "openid profile email"
  federation_authorization_url        = "http://localhost:8800/auth"
  federation_token_url                = "http://localhost:8800/token"
  federation_user_info_url            = "http://localhost:8800/userinfo"
  federation_jwk_set_url              = "http://localhost:8800/auth/jwks"
  federation_auth_script_text         = local.auth_script_1
  federation_user_mapping_script_text = local.user_mapping_script_1
  token_introspection_client_id       = "smile"
  token_introspection_client_secret   = "client_secret_goes_here"
  validation_jwk_text                 = ""
  validation_jwk_file                 = ""
}

resource "smilecdr_openid_identity_provider" "example_idp2" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  issuer                              = "https://example2.com/auth/issuer"
  name                                = "example_idp2"
  federation_authorization_url        = "http://example2.com/auth/authorize"
  federation_request_scopes           = "openid profile fhirUser"
  federation_token_url                = "https://example2.com/auth/token"
  federation_user_info_url            = "https://example2.com/auth/userinfo"
  federation_jwk_set_url              = "https://example2.com/auth/jwks"
  federation_auth_script_text         = local.auth_script_2
  federation_user_mapping_script_text = local.user_mapping_script_2
  token_introspection_client_id       = "smile"
  token_introspection_client_secret   = "client_secret_goes_here"
  validation_jwk_text                 = ""
  validation_jwk_file                 = ""
}