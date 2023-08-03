# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_openid_identity_provider" "example_idp1" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  issuer                              = "https://example1.com/auth/issuer"
  name                                = "example_idp"
  federation_authorization_url        = "http://localhost:8800/auth"
  federation_request_scopes           = "openid profile email"
  federation_token_url                = "http://localhost:8800/token"
  federation_user_info_url            = "http://localhost:8800/userinfo"
  federation_jwk_set_url              = "http://localhost:8800/auth/jwks"
  federation_auth_script_text         = local.authScriptFile1
  federation_user_mapping_script_text = local.userMappingScriptFile
  token_introspection_client_id       = "smile"
  token_introspection_client_secret   = "client_secret"
  validation_jwk_text                 = ""
  validation_jwk_file                 = ""
}

resource "smilecdr_openid_identity_provider" "example_idp2" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  issuer                              = "https://example2.com/auth/issuer"
  name                                = "idp5"
  federation_authorization_url        = "http://example2.com/auth/authorize"
  federation_request_scopes           = "openid profile fhirUser"
  federation_token_url                = "https://example2.com/auth/token"
  federation_user_info_url            = "https://example2.com/auth/userinfo"
  federation_jwk_set_url              = "https://example2.com/auth/jwks"
  federation_auth_script_text         = local.authScriptFile2
  federation_user_mapping_script_text = local.userMappingScriptFile
  token_introspection_client_id       = "smile"
  token_introspection_client_secret   = "client_secret"
  validation_jwk_text                 = ""
}
