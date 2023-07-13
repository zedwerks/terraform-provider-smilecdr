# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

terraform {
  required_providers {
    smilecdr = {
      source  = "local.providers/zedwerks/smilecdr"
      version = "0.2.0"
    }
  }
}

provider "smilecdr" {
  base_url = "http://localhost:9000"
  username = "admin"
  password = "password"
}


resource "smilecdr_openid_client" "client1" {

  node_id                       = "Master"
  module_id                     = "smart_auth"
  access_token_validity_seconds = 300
  allowed_grant_types = [
    "authorization_code",
    "refresh_token",
    "client_credentials",
    "password",
    "implicit"
  ]
  auto_approve_scopes = [
    "openid",
    "profile",
    "email",
    "patient/*.*",
    "launch",
    "launch/patient",
    "launch/encounter",
    "offline_access"
  ]
  auto_grant_scopes = [
    "openid",
    "profile",
    "email",
    "patient/*.*",
    "launch",
    "launch/patient",
    "launch/encounter",
    "offline_access"
  ]
  client_id                      = "client1"
  client_name                    = "Client 1"
  fixed_scope                    = false
  refresh_token_validity_seconds = 86400
  registered_redirect_uris       = ["http://localhost:3000"]
  scopes = [
    "openid",
    "profile",
    "email",
    "patient/*.*",
    "launch",
    "launch/patient",
    "launch/encounter",
    "offline_access"
  ]
  secret_required           = false
  secret_client_can_change  = false
  enabled                   = true
  can_introspect_any_tokens = false
  can_introspect_own_tokens = false
  always_require_approval   = false
  can_reissue_tokens        = false
  remember_approved_scopes  = false
  attestation_accepted      = false
  public_jwks_uri           = ""
  created_by_app_sphere     = false

}
