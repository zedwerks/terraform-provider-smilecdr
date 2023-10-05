# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0
#
# This is an example of how to use the smart inbound security module.
# An inbound security module validates tokens only. It is not an authorization server
# and does not issue tokens.  It is used to validate tokens issued by an external
# authorization server.

resource "smilecdr_smart_inbound_security" "ex1_smart_inbound" {
  module_id                                       = "ex1_smart_inbound"
  node_id                                         = "Master"
  enforce_approved_scopes_to_restrict_permissions = false
  trust_intra_cluster_tokens_modules              = ""
  cache_authentication_seconds                    = 300
  key_validation_prevent_token_key_reuse          = false
  key_validation_require_key_expiry               = false
  smart_configuration_scopes_supported            = "openid fhirUser patient/*.read"
  token_endpoint                                  = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/token"
  authorization_endpoint                          = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/auth"
  management_endpoint                             = "http://localhost:8080/auth/realms/poc/account"
  introspection_endpoint                          = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/token/introspect"
  revocation_endpoint                             = "http://localhost:8080/auth/realms/poc/protocol/openid-connect/revoke"
  introspection_client_jwks_cache_mins            = 10
  introspection_client_truststore_file            = "file://somefile.jks"
  callback_script_text                            = local.inbound_callback_script
  tfa_totp_issuer_name                            = "issuerName"
  tfa_totp_lock_after_failed_attempts             = 2
  seed_servers_file                               = "seedServers.txt"
  debug_enabled                                   = false
  debug_host_address                              = "localhost"
  debug_secure                                    = false
  debug_suspend                                   = false
  debug_port                                      = 5005
  debug_path                                      = "/debug"
}
