# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0


resource "smilecdr_openid_identity_provider" "example_idp1" {
  node_id                             = "Master"
  module_id                           = "smart_auth"
  name                                = "idp1"
  issuer                              = "https://example.com/idp1"
  token_introspection_client_id       = "https://example.com/idp1/jwks"
  token_introspection_client_secret   = "https://example.com/idp1/auth"
  validation_jwk_text                 = ""
  validation_jwk_file                 = ""
  federation_registration_id          = ""
  federation_request_scopes           = ""
  federation_authorization_url        = ""
  federation_token_url                = ""
  federation_user_info_url            = ""
  federation_auth_script_text         = ""
  federation_auth_script_file         = ""
  federation_user_mapping_script_text = ""
  fhir_endpoint_url                   = ""
  auth_well_known_config_url          = ""
  notes                               = ""
  custom_token_params                 = ""
  response_type                       = ""
  organization_id                     = ""
  audience                            = ""
  #archived_at = ""
}
