# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0
#
# This is an example of how to use the smart inbound security module.
# An inbound security module validates tokens only. It is not an authorization server
# and does not issue tokens.  It is used to validate tokens issued by an external
# authorization server.



resource "smilecdr_smart_inbound_security" "example_inbound" {
  module_id = "example_smart_inbound"
  node_id   = "Master"

  cache_authentication_seconds                    = 300
  key_validation_prevent_token_key_reuse          = true
  key_validation_require_key_expiry               = true
  enforce_approved_scopes_to_restrict_permissions = true
  smart_configuration_scopes_supported            = "openid fhirUser patient/*.read"
  token_endpoint                                  = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/token"
  authorization_endpoint                          = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/auth"
  introspection_endpoint                          = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/token/introspect"
  revocation_endpoint                             = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/revoke"
}
