# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_smart_outbound_security" "smart_auth_outbound" {
  module_id                                  = "smart_auth_outbound_min"
  node_id                                    = "Master"
  http_listener_context_path                 = "/"
  http_listener_unhealthy_response_code      = 503
  http_listener_port                         = 9202
  oidc_issuer_url                      = "http://keycloak:8080/auth/realms/poc"
  oidc_smart_capabilities_list = [
    "launch-ehr",
    "launch-standalone",
    "launch-patient",
    "launch",
    "patient/*.read",
    "patient/*.write",
    "patient/*.*"
  ]
  oidc_federate_mode_enabled = true
  smart_callback_post_authorize_script_text = local.post_authorize_callback
  codap_enabled                   = false
  smart_authorization_email_from_address      = "noreply@phsa.ca"
  smart_authorization_enforce_approved_scopes = true
  smart_authorization_scopes_supported = "launch, fhirUser, openid, profile, patient/*.read, patient/*.write, patient/*.*"
  sessions_in_memory = false
  sessions_max_concurrent_sessions_per_user = 3
  dependency_local_inbound_security = "SECURITY_IN_UP"
  dependency_fhir_persistence_module = "PERSISTANCE_ALL"
  dependency_saml_authentication_module = "SECURITY_IN_SAML"
}
