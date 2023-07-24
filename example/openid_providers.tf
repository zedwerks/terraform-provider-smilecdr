# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

locals {
  smilecdr_version = "2021.05.R01"
  authScriptFile = file("${path.module}/js/federationAuthScript.js")
  userMappingScriptFile = file("${path.module}/js/federationUserMappingScript.js")
}

resource "smilecdr_openid_identity_provider" "idp4" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  issuer                              = "https://idp4.com/auth/issuer"
  name                                = "idp4"
  federation_authorization_url        = "http://localhost:8800/auth"
  federation_request_scopes           = "openid profile email"
  federation_token_url                = "http://localhost:8800/token"
  federation_user_info_url            = "http://localhost:8800/userinfo"
  federation_jwk_set_url              = "//localhost:8800/auth/jwks"
  federation_auth_script_text         = "print('hello world')"
  federation_user_mapping_script_text = "// user mapping script goes here"
  token_introspection_client_id       = "smile "
  token_introspection_client_secret   = "client_secret"
  validation_jwk_text                 = ""
  validation_jwk_file                 = ""
}

resource "smilecdr_openid_identity_provider" "idp5" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  issuer                              = "https://idp5.com/auth/issuer"
  name                                = "idp5"
  federation_authorization_url        = "htt://idp5.com/auth/authorize"
  federation_request_scopes           = "openid profile fhirUser"
  federation_token_url                = "https://idp5.com/auth/token"
  federation_user_info_url            = "https://idp5.com/auth/userinfo"
  federation_jwk_set_url              = "https://idp5.com/auth/jwks"
  federation_auth_script_text         = local.authScriptFile
  federation_user_mapping_script_text = local.userMappingScriptFile
  token_introspection_client_id       = "smile"
  token_introspection_client_secret   = "client_secret"
}
