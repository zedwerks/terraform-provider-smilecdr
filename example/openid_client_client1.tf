# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_openid_client" "client1" {
  node_id                       = "Master"
  module_id                     = "smart_auth"
  access_token_validity_seconds = 300
  allowed_grant_types           = ["REFRESH_TOKEN", "CLIENT_CREDENTIALS"]
  auto_approve_scopes = ["openid", "profile",
    "fhirUser",
    "launch",
    "launch/patient",
    "offline_access"
  ]
  auto_grant_scopes = ["openid", "offline_access"]
  client_id         = "client1"
  client_name       = "Client1"
  client_secrets {
    secret     = "secret1234569900"
    activation = "2023-08-05T00:09:53.702+00:00"
  }
  client_secrets {
    secret = "secret23456789aabb"
  }
  fixed_scope                    = false
  refresh_token_validity_seconds = 86400
  registered_redirect_uris       = ["http://example-client1.com:6000", "http://example-client1.com:6001"]
  scopes = [
    "openid",
    "profile",
    "fhirUser",
    "patient/*.read",
    "launch",
    "launch/patient",
    "offline_access"
  ]
  secret_required           = true
  secret_client_can_change  = false
  enabled                   = true
  can_introspect_any_tokens = true
  can_introspect_own_tokens = false
  always_require_approval   = false
  can_reissue_tokens        = false
  remember_approved_scopes  = false
  attestation_accepted      = false
  jwks_url                  = "http://example-client1.com/jwks"
  permissions {
    permission = "FHIR_WRITE_ALL_IN_COMPARTMENT"
    argument   = "Patient/123"
  }
  permissions {
    permission = "ROLE_FHIR_CLIENT"
  }
}
