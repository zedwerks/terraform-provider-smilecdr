# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0


resource "smilecdr_openid_client" "postman" {
  depends_on = [ smilecdr_smart_outbound_security.smart_auth_federated ]
  node_id                       = "Master"
  module_id                     = "smart_auth_federated"
  access_token_validity_seconds = 300
  allowed_grant_types           = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
  auto_approve_scopes = ["openid", "profile",
    "fhirUser",
    "launch",
    "launch/patient",
    "offline_access",
  ]
  auto_grant_scopes = ["openid", "offline_access"]
  client_id         = "postman"
  client_name       = "Postman on localhost"
  fixed_scope                    = false
  refresh_token_validity_seconds = 86400
  registered_redirect_uris       = ["https://oauth.pstmn.io/v1/callback" ]
  scopes = [
    "openid",
    "profile",
    "fhirUser",
    "patient/*.read",
    "launch",
    "launch/patient",
    "offline_access"
  ]
  secret_required           = false
  secret_client_can_change  = false
  enabled                   = true
  can_introspect_any_tokens = true
  can_introspect_own_tokens = false
  always_require_approval   = false
  can_reissue_tokens        = false
  remember_approved_scopes  = true
  attestation_accepted      = false
  permissions {
    permission = "ROLE_FHIR_CLIENT"
  }
}
