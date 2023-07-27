# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0


resource "smilecdr_openid_client" "client_lra_demo" {
  node_id                       = "Master"
  module_id                     = "smart_auth"
  access_token_validity_seconds = 300
  allowed_grant_types = [
    "AUTHORIZATION_CODE",
  "REFRESH_TOKEN", ]
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
  client_id                      = "client_lra_demo"
  client_name                    = "Client LRA Demo"
  fixed_scope                    = false
  refresh_token_validity_seconds = 86400
  registered_redirect_uris       = ["http://localhost:6000"]
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
}

