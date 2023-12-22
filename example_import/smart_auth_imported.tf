

## Imported resource "smart_auth"
resource "smilecdr_smart_outbound_security" "smart_auth" {
  module_id                            = "smart_auth"
  node_id                              = "Master"
  http_listener_port                   = 9200
  oidc_issuer_url                      = "http://smilecdr:9200"
  smart_authorization_scopes_supported = "fhirUser launch/patient patient/*.read user/*.read offline_access openid profile email"
  cors_enabled                         = true
  oidc_federate_mode_enabled           = true
  jwks_keystore_id                     = "default-keystore"
  oidc_smart_capabilities_list = [
    "client-public",
    "client-confidential-symmetric",
    "sso-openid-connect",
    "context-standalone-patient",
    "permission-patient",
    "permission-user",
  ]
  oidc_pkce_required                        = true
  davinci_native_consent_handling           = false
  smart_callback_post_authorize_script_text = local.post_authorize_callback
  dependency_local_inbound_security         = "SECURITY_IN_UP"
}
