// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func dataSourceSmartOutboundSecurity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSmartOutboundSecurityRead,
		Schema: map[string]*schema.Schema{
			"module_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique module ID of the module to be configured.",
			},
			"module_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The module type of the module to be configured.",
			},
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The node ID of the node to be configured.",
			},
			// User Authentication Options ------------------------
			"anonymous_account_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username to use for the anonymous user account. This account will be used for anonymous requests (i.e. requests without credentials).",
			},
			"anonymous_access_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If enabled, anonymous requests (i.e. requests without credentials) will be allowed to proceed. This means that they will not be blocked by the security manager, and they will instead proceed under the authority of the designated anonymous user. Only roles and permissions that have been assigned to the anonymous user will be granted to these requests. See Anonymous Access for more information.",
			},
			"saml_authentication_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If enabled, the server will allow authentication via SAML. This will enable the SAML authentication module, which will allow users to authenticate via SAML. See SAML Authentication for more information.",
			},
			// CORS Options ------------------------
			"cors_allowed_headers": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comma-separated list of allowable request headers for the CORS filter. These will be added in addition to the default headers required for Smile CDR's default functionality.",
			},
			"cors_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Should this endpoint allow the use of CORS? Enable this item only if you understand what it is doing.",
			},
			"cors_origins": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A comma-separated list of allowable origins for the CORS filter. For example: https://example.com, https://try.smilecdr.com:9201. You may also use the wildcard value * to allow CORS for all domains, however this is generally not considered a good practice for production systems serving sensitive data.",
			},
			// Davinci Options ------------------------
			"davinci_native_consent_handling": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Activates handling of consents via the persistence module. See Da Vinci Health Record Exchange.",
			},
			// HTTP Access Log Options ------------------------
			"http_access_log_appenders": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of appenders to use for HTTP access logging. Each appender should be specified as a single line in the format: appender-name",
			},
			// HTTP Listener Options ------------------------
			"http_listener_bind_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_listener_context_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_listener_endpoint_health_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_listener_unhealthy_response_code": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_listener_https_forwarding_assumed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"http_listener_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_listener_respect_forward_headers": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// HTTP Request Pool Options ------------------------
			"http_request_maximum_request_header_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_request_maximum_response_header_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_request_read_idle_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_request_thread_pool_accept_queue_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_request_thread_pool_max_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_request_thread_pool_min_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// HTTP Security Options --------------------------------
			"http_security_block_http_head": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set, the server will reject the HTTP HEAD verb. This verb is considered insecure in some environments.",
			},
			"http_security_block_http_options": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set, the server will reject the HTTP OPTIONS verb. This verb is considered insecure in some environments.",
			},
			"http_security_custom_response_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Custom headers to add to all responses. Each header should be specified as a single line in the format: Header-Name: Header-Value",
			},
			"http_security_frame_options_allow_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This setting can be used to set the X-Frame-Options header. Leave this setting blank (the default) in order to set a value of DENY. See Frame Options for more information.",
			},
			"http_security_pin_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set, the server will always use the given host name instead of respecting the Host header. Comma-separated list of host names to pin to.",
			},
			"http_security_suppress_error_details": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"http_security_suppress_platform_info": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// JavaScript Execution Environment Options ------------------------
			"javascript_debug_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable remote JavaScript debugging.",
			},
			"javascript_debug_host_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hostname of the server running Smile CDR",
			},
			"javascript_debug_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"javascript_debug_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"javascript_debug_secure": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"javascript_debug_suspend": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// JWKS Options --------------------------------
			"jwks_keystore_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the ID of the keystore to use. The keystore defines the signing keys and can be managed in admin console. This config overrides all other configs in this section.",
			},
			// OIDC Token Validation Options ------------------------
			"oidc_http_client_jwks_cache_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"oidc_http_client_truststore_file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oidc_http_client_truststore_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// OpenID Connect (OIDC) Options ------------------------
			"oidc_pkce_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If this setting is enabled, the server will require the use of PKCE for all Authorization Code SMART Auth flows. Enabling this setting also disallows the use of the OAuth2 Implicit Grant type, since this flow does not support PKCE.",
			},
			"oidc_pkce_plain_challenge_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If this setting is enabled, the server will allow the use of the plain PKCE challenge method. This is not recommended, but is supported for backwards compatibility.",
			},
			"oidc_cache_authorization_tokens": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"oidc_client_secret_encoding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Select the hashing algorithm to use when storing client secrets. Note that the value selected here will apply only to newly created secrets, and this may be changed at any time without affecting existing secrets.",
			},
			"oidc_client_secret_expiry_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Select the expiry duration in days for Smile CDR generated client secrets. Note this value will be added to the activation date of the secret to calculate the expiration date for the secret during the client creation process via the REST path register-client-and-generate-secret.",
			},
			"oidc_issuer_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the URL that will be placed in OpenID Connect tokens as the iss (issuer) token. The value should be the URL to the identity server.",
			},

			"oidc_rotate_token_after_use": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If enabled, each time a refresh token is used to obtain a new access token, the refresh token will be invalidated and a new one automatically issued with the new access token.",
			},
			"oidc_smart_capabilities_list": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "A list of SMART capabilities to advertise in the .well-known/smart-configuration.",
			},
			// OAuth2/OIDC Federation Options ------------------------
			"oidc_federate_mode_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When enabled, this server will federate to a federated OAuth2/OIDC server instead of prompting the user for credentials. See Federated OAuth2/OIDC Login for more information.",
			},
			// SMART Callback Script Options ------------------------
			"smart_callback_post_authorize_script_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If supplied, provides a script that will be invoked after various authorization flows complete. See SMART Callback Script for details on how this script works. Values should be prefixed with file: or classpath:.",
			},
			"smart_callback_post_authorize_script_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If supplied, provides a script that will be invoked after various authorization flows complete. See SMART Callback Script for details on how this script works. Values should be prefixed with file: or classpath:.",
			},
			// CODAP Options ------------------------
			"codap_authorization_script_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When using CODAP, a callback script must be provided. This script is used to assess the incoming identity assertion and provide the appropriate authorization. Values should be prefixed with file: or classpath:.",
			},
			"codap_authorization_script_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When using CODAP, a callback script must be provided. This script is used to assess the incoming identity assertion and provide the appropriate authorization.",
			},
			"codap_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// SMART Login Skin Options ------------------------
			"smart_login_skin_approval_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the approval page, e.g. /userapprove.html",
			},
			"smart_login_skin_context_selection_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the Context Selection page, e.g. /userselect_context.html",
			},
			"smart_login_skin_error_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the error page, e.g. /usererror.html",
			},
			"smart_login_skin_user_registration_forgot_password_template_step1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the first page of the user self registration forgot password flow",
			},
			"smart_login_skin_user_registration_forgot_password_template_step2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the second page of the user self registration forgot password flow",
			},
			"smart_login_skin_user_registration_forgot_password_template_step3": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the third page of the user self registration forgot password flow",
			},
			"smart_login_skin_federated_oath2_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the federated OAuth2/OIDC Login Page, e.g. /userlogin_oauth2.html. This setting is not used unless Federated OAuth2/OIDC Login is enabled, and may be left blank otherwise.",
			},
			"smart_login_skin_login_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the interactive login page, e.g. /userlogin.html",
			},
			"smart_login_skin_user_registration_template_step1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the first page of the user self registration flow, e.g. /userregister_step1.html",
			},
			"smart_login_skin_user_registration_template_step2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the second page of the user self registration flow, e.g. /userregister_step2.html",
			},
			"smart_login_skin_session_management_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the session management page, e.g. /sessionmanagement.html",
			},

			"smart_login_skin_2fa_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the two factor authentication code entry page, e.g. /usertfa.html",
			},
			"smart_login_skin_terms_of_service_template": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the path within the WebJar for the click-wrap terms-of-service agreement entry page, e.g. /usertos.html",
			},
			"smart_login_skin_webjar_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This is the ID of the WebJar to use as a skin for the SMART Outbound Security module for login and approval screens. This should take the form groupId:artifactId:versionId.",
			},
			// SMART Login Terms of Service Options ------------------------
			"smart_login_terms_of_service_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// SMART Authorization Options ------------------------
			"smart_authorization_allowed_audience_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Space-separated list of allowed resource URLs as the 'audience' parameter during authentication flow. If left empty, no validation is performed.",
			},
			"smart_authorization_email_from_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Forgotten password related emails will be sent from this email address.",
			},
			"smart_authorization_enforce_approved_scopes": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When enabled, permission will be stripped from a user's session if they are not supported by an approved SMART on FHIR scope. For example, any FHIR write permissions will be removed from a session if the user has not approved (or a client is set to auto-approve) a scope such as Patient/*.write",
			},
			"smart_authorization_scopes_supported": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A space separated list of scopes to advertise as supported in the .well-known/smart-configuration.",
			},

			// SMART Definitions Seeding Options ------------------------
			"openid_connect_server_pre_seed_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provides the location of a file to use to pre-seed OpenID Connect Client definitions at startup time. See Pre-Seeding for more information",
			},
			"openid_connect_client_pre_seed_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provides the location of a file to use to pre-seed OpenID Connect Server definitions at startup time. See Pre-Seeding for more information.",
			},
			// Sessions Options ------------------------
			"sessions_in_memory": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If enabled, any HTTP sessions created for this listener will be stored only in memory, as opposed to being persisted in the database. This may lead to a performance boost in some situations but also prevents sessions from working in some clustered configurations or surviving a restart of the system. Note that not all listeners even create sessions (e.g. FHIR endpoints do not) so this setting may have no effect",
			},
			"sessions_max_concurrent_sessions_per_user": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "If set to a value greater than zero, this setting will limit the number of concurrent sessions that a single user can have. If a user attempts to create a new session when they already have the maximum number of sessions, the oldest session will be terminated. This setting is useful for preventing users from sharing their credentials with others.",
			},
			"sessions_scavenger_interval_ms": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of milliseconds between session scavenger passes.",
			},
			"sessions_timeout_mins": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of minutes that a user session can sit idle before it is eligible to expire.",
			},
			// TLS Options ------------------------
			"tls_client_auth_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Should the listener for this module require incoming connections to authenticate using TLS Client Authentication?",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Should the listener for this module require TLS (i.e. SSL or HTTPS) encryption for incoming connections?",
			},
			"tls_keystore_filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The filename for the TLS KeyStore used to hold private keys for TLS connections. This can be in the format classpath:path/to/file.p12 or file:///path/to/file.p12. Valid file extensions are .jks (Java Keystore) or .p12 (PKCS#12 store).",
			},
			"tls_keystore_key_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alias for the specific key within the KeyStore that should be selected for incoming TLS connections.",
			},
			"tls_keystore_key_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password for the specific key within the KeyStore (leave blank if the key has no password).",
			},
			"tls_keystore_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password for the TLS KeyStore (leave blank if the store has no password).",
			},
			"tls_cipher_allow_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "f specified, contains a space-separated list of ciphers that are permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information.",
			},
			"tls_cipher_deny_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If specified, contains a space-separated list of ciphers that are not permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information",
			},
			"tls_protocol_allow_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If specified, contains a space-separated list of protocols that are permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information.",
			},
			"tls_protocol_deny_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If specified, contains a space-separated list of protocols that are not permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information.",
			},

			"tls_truststore_filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The filename for the TLS TrustStore used to hold trusted certificates for TLS connections. This can be in the format classpath:path/to/file.p12 or file:///path/to/file.p12. Valid file extensions are .jks (Java Keystore) or .p12 (PKCS#12 store).",
			},
			"tls_truststore_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password for the TLS TrustStore (leave blank if the store has no password).",
			},
			// Dependency Options ------------------------
			"dependency_local_inbound_security": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The inbound security module to use for authenticating and authorizing users to this module where authentication requires a username and password.",
			},
			"dependency_fhir_persistence_module": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FHIR Storage module to associate with this module.",
			},
			"dependency_saml_authentication_module": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SAML Inbound Security module to use when performing a SAML user authentication.",
			},
			"dependency_self_registration_provider_module": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This can be supplied to some interactive modules in order to support self-registration of users.",
			},
		},
	}
}

func dataSourceSmartOutboundSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Get the module ID from the resource data
	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)

	c := m.(*smilecdr.Client)

	moduleConfig, err := c.GetModuleConfig(nodeId, moduleId)

	// map from moduleConfig to resourceData
	if err != nil {
		return diag.FromErr(err)
	}

	err = setSmartOutboundSecurityData(ctx, d, &moduleConfig)

	return diag.FromErr(err)
}
