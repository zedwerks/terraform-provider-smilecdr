# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

resource "smilecdr_smart_outbound_security" "smart_auth_federated" {
		module_id                             = "test_auth_outbound3"
		node_id                               = "Master"
		http_listener_context_path            = "/auth"
		http_listener_unhealthy_response_code = 503
		http_listener_port                    = 9201
		oidc_issuer_url                       = "http://keycloak:8080/auth/realms/poc"
		oidc_federate_mode_enabled            = true
		smart_authorization_enforce_approved_scopes = true
		smart_authorization_scopes_supported        = "launch fhirUser openid profile patient/*.read"
		sessions_max_concurrent_sessions_per_user   = 3
		dependency_fhir_persistence_module          = "PERSISTENCE_ALL"
}
