// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zedwerks/terraform-smilecdr/provider/helper/validations"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func resourceSmartOutboundSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmartOutboundSecurityCreate,
		ReadContext:   resourceSmartOutboundSecurityRead,
		UpdateContext: resourceSmartOutboundSecurityUpdate,
		DeleteContext: resourceSmartOutboundSecurityDelete,
		Schema: map[string]*schema.Schema{
			"module_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique module ID of the module to be configured.",
			},
			"module_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The module type of the module to be configured.",
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Master",
				Description: "The node ID of the node to be configured.",
			},
			// User Authentication Options ------------------------
			"anonymous_account_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ANONYMOUS",
				Description: "The username to use for the anonymous user account. This account will be used for anonymous requests (i.e. requests without credentials).",
			},
			"anonymous_access_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, anonymous requests (i.e. requests without credentials) will be allowed to proceed. This means that they will not be blocked by the security manager, and they will instead proceed under the authority of the designated anonymous user. Only roles and permissions that have been assigned to the anonymous user will be granted to these requests. See Anonymous Access for more information.",
			},
			"saml_authentication_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, the server will allow authentication via SAML. This will enable the SAML authentication module, which will allow users to authenticate via SAML. See SAML Authentication for more information.",
			},
			// CORS Options ------------------------
			"cors_allowed_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A comma-separated list of allowable request headers for the CORS filter. These will be added in addition to the default headers required for Smile CDR's default functionality.",
			},
			"cors_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should this endpoint allow the use of CORS? Enable this item only if you understand what it is doing.",
			},
			"cors_origins": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				Description: "A comma-separated list of allowable origins for the CORS filter. For example: https://example.com, https://try.smilecdr.com:9201. You may also use the wildcard value * to allow CORS for all domains, however this is generally not considered a good practice for production systems serving sensitive data.",
			},
			// Davinci Options ------------------------
			"davinci_native_consent_handling": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Activates handling of consents via the persistence module. See Da Vinci Health Record Exchange.",
			},
			// HTTP Access Log Options ------------------------
			"http_access_log_appenders": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of appenders to use for HTTP access logging. Each appender should be specified as a single line in the format: appender-name",
			},
			// HTTP Listener Options ------------------------
			"http_listener_bind_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0",
			},
			"http_listener_context_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},
			"http_listener_endpoint_health_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/endpoint-health",
			},
			"http_listener_unhealthy_response_code": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  503, // 503 Service Unavailable
			},
			"http_listener_https_forwarding_assumed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http_listener_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"http_listener_respect_forward_headers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// HTTP Request Pool Options ------------------------
			"http_request_maximum_request_header_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
			},
			"http_request_maximum_response_header_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8,
			},
			"http_request_read_idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30000,
			},
			"http_request_thread_pool_accept_queue_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"http_request_thread_pool_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"http_request_thread_pool_min_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			// HTTP Security Options --------------------------------
			"http_security_block_http_head": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set, the server will reject the HTTP HEAD verb. This verb is considered insecure in some environments.",
			},
			"http_security_block_http_options": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set, the server will reject the HTTP OPTIONS verb. This verb is considered insecure in some environments.",
			},
			"http_security_custom_response_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Custom headers to add to all responses. Each header should be specified as a single line in the format: Header-Name: Header-Value",
			},
			"http_security_frame_options_allow_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This setting can be used to set the X-Frame-Options header. Leave this setting blank (the default) in order to set a value of DENY. See Frame Options for more information.",
			},
			"http_security_pin_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If set, the server will always use the given host name instead of respecting the Host header. Comma-separated list of host names to pin to.",
			},
			"http_security_suppress_error_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http_security_suppress_platform_info": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// JavaScript Execution Environment Options ------------------------
			"javascript_debug_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable remote JavaScript debugging.",
			},
			"javascript_debug_host_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "localhost",
				Description: "The hostname of the server running Smile CDR",
			},
			"javascript_debug_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"javascript_debug_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  9930,
			},
			"javascript_debug_secure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"javascript_debug_suspend": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// JWKS Options --------------------------------
			"jwks_keystore_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the ID of the keystore to use. The keystore defines the signing keys and can be managed in admin console. This config overrides all other configs in this section.",
			},
			// OIDC Token Validation Options ------------------------
			"oidc_http_client_jwks_cache_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"oidc_http_client_truststore_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_http_client_truststore_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// OpenID Connect (OIDC) Options ------------------------
			"oidc_cache_authorization_tokens": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3000,
			},
			"oidc_client_secret_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "BCRYPT_12_ROUND",
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringInSlice([]string{
					"SHA256_1000_ROUND",
					"SHA256_10000_ROUND",
					"SHA256_100000_ROUND",
					"PBKDF2_256_1000_RND",
					"PBKDF2_256_10000_RND",
					"PBKDF2_256_100000_RND",
					"BCRYPT_10_ROUND",
					"BCRYPT_12_ROUND",
					"BCRYPT_14_ROUND",
					"BCRYPT_16_ROUND"}, false)),
				Description: "Select the hashing algorithm to use when storing client secrets. Note that the value selected here will apply only to newly created secrets, and this may be changed at any time without affecting existing secrets.",
			},
			"oidc_client_secret_expiry_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     365,
				Description: "Select the expiry duration in days for Smile CDR generated client secrets. Note this value will be added to the activation date of the secret to calculate the expiration date for the secret during the client creation process via the REST path register-client-and-generate-secret.",
			},
			"oidc_issuer_url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				Description:      "This is the URL that will be placed in OpenID Connect tokens as the iss (issuer) token. The value should be the URL to the identity server.",
			},
			"oidc_pkce_plain_challenge_supported": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If this setting is enabled, the server will allow the use of the plain PKCE challenge method. This is not recommended, but is supported for backwards compatibility.",
			},
			"oidc_pkce_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If this setting is enabled, the server will require the use of PKCE for all Authorization Code SMART Auth flows. Enabling this setting also disallows the use of the OAuth2 Implicit Grant type, since this flow does not support PKCE.",
			},
			"oidc_rotate_token_after_use": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, each time a refresh token is used to obtain a new access token, the refresh token will be invalidated and a new one automatically issued with the new access token.",
			},
			"oidc_smart_capabilities_list": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "A list of SMART capabilities to advertise in the .well-known/smart-configuration.",
			},
			// OAuth2/OIDC Federation Options ------------------------
			"oidc_federate_mode_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When enabled, this server will federate to a federated OAuth2/OIDC server instead of prompting the user for credentials. See Federated OAuth2/OIDC Login for more information.",
			},
			// SMART Callback Script Options ------------------------
			"smart_callback_post_authorize_script_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If supplied, provides a script that will be invoked after various authorization flows complete. See SMART Callback Script for details on how this script works. Values should be prefixed with file: or classpath:.",
			},
			"smart_callback_post_authorize_script_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If supplied, provides a script that will be invoked after various authorization flows complete. See SMART Callback Script for details on how this script works. Values should be prefixed with file: or classpath:.",
			},
			// CODAP Options ------------------------
			"codap_authorization_script_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When using CODAP, a callback script must be provided. This script is used to assess the incoming identity assertion and provide the appropriate authorization. Values should be prefixed with file: or classpath:.",
			},
			"codap_authorization_script_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "When using CODAP, a callback script must be provided. This script is used to assess the incoming identity assertion and provide the appropriate authorization.",
			},
			"codap_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// SMART Login Skin Options ------------------------
			"smart_login_skin_approval_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the approval page, e.g. /userapprove.html",
			},
			"smart_login_skin_context_selection_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the Context Selection page, e.g. /userselect_context.html",
			},
			"smart_login_skin_error_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the error page, e.g. /usererror.html",
			},
			"smart_login_skin_user_registration_forgot_password_template_step1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the first page of the user self registration forgot password flow",
			},
			"smart_login_skin_user_registration_forgot_password_template_step2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the second page of the user self registration forgot password flow",
			},
			"smart_login_skin_user_registration_forgot_password_template_step3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the third page of the user self registration forgot password flow",
			},
			"smart_login_skin_federated_oath2_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the federated OAuth2/OIDC Login Page, e.g. /userlogin_oauth2.html. This setting is not used unless Federated OAuth2/OIDC Login is enabled, and may be left blank otherwise.",
			},
			"smart_login_skin_login_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the interactive login page, e.g. /userlogin.html",
			},
			"smart_login_skin_user_registration_template_step1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the first page of the user self registration flow, e.g. /userregister_step1.html",
			},
			"smart_login_skin_user_registration_template_step2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the second page of the user self registration flow, e.g. /userregister_step2.html",
			},
			"smart_login_skin_session_management_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the session management page, e.g. /sessionmanagement.html",
			},

			"smart_login_skin_2fa_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the two factor authentication code entry page, e.g. /usertfa.html",
			},
			"smart_login_skin_terms_of_service_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the path within the WebJar for the click-wrap terms-of-service agreement entry page, e.g. /usertos.html",
			},
			"smart_login_skin_webjar_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the ID of the WebJar to use as a skin for the SMART Outbound Security module for login and approval screens. This should take the form groupId:artifactId:versionId.",
			},
			// SMART Login Terms of Service Options ------------------------
			"smart_login_terms_of_service": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// SMART Authorization Options ------------------------
			"smart_authorization_allowed_audience_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Space-separated list of allowed resource URLs as the 'audience' parameter during authentication flow. If left empty, no validation is performed.",
			},
			"smart_authorization_email_from_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "noreply@unknown.com",
				Description: "Forgotten password related emails will be sent from this email address.",
			},
			"smart_authorization_enforce_approved_scopes": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "When enabled, permission will be stripped from a user's session if they are not supported by an approved SMART on FHIR scope. For example, any FHIR write permissions will be removed from a session if the user has not approved (or a client is set to auto-approve) a scope such as Patient/*.write",
			},
			"smart_authorization_scopes_supported": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A space separated list of scopes to advertise as supported in the .well-known/smart-configuration.",
			},

			// SMART Definitions Seeding Options ------------------------
			"openid_connect_server_pre_seed_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Provides the location of a file to use to pre-seed OpenID Connect Client definitions at startup time. See Pre-Seeding for more information",
			},
			"openid_connect_client_pre_seed_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Provides the location of a file to use to pre-seed OpenID Connect Server definitions at startup time. See Pre-Seeding for more information.",
			},
			// Sessions Options ------------------------
			"sessions_in_memory": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, any HTTP sessions created for this listener will be stored only in memory, as opposed to being persisted in the database. This may lead to a performance boost in some situations but also prevents sessions from working in some clustered configurations or surviving a restart of the system. Note that not all listeners even create sessions (e.g. FHIR endpoints do not) so this setting may have no effect",
			},
			"sessions_max_concurrent_sessions_per_user": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "If set to a value greater than zero, this setting will limit the number of concurrent sessions that a single user can have. If a user attempts to create a new session when they already have the maximum number of sessions, the oldest session will be terminated. This setting is useful for preventing users from sharing their credentials with others.",
			},
			"sessions_scavenger_interval_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60000,
				Description: "The number of milliseconds between session scavenger passes.",
			},
			"sessions_timeout_mins": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "The number of minutes that a user session can sit idle before it is eligible to expire.",
			},
			// TLS Options ------------------------
			"tls_client_auth_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should the listener for this module require incoming connections to authenticate using TLS Client Authentication?",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should the listener for this module require TLS (i.e. SSL or HTTPS) encryption for incoming connections?",
			},
			"tls_keystore_filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The filename for the TLS KeyStore used to hold private keys for TLS connections. This can be in the format classpath:path/to/file.p12 or file:///path/to/file.p12. Valid file extensions are .jks (Java Keystore) or .p12 (PKCS#12 store).",
			},
			"tls_keystore_key_alias": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "The alias for the specific key within the KeyStore that should be selected for incoming TLS connections.",
			},
			"tls_keystore_key_password": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "The password for the specific key within the KeyStore (leave blank if the key has no password).",
			},
			"tls_keystore_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for the TLS KeyStore (leave blank if the store has no password).",
			},
			"tls_cipher_allow_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "f specified, contains a space-separated list of ciphers that are permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information.",
			},
			"tls_cipher_deny_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If specified, contains a space-separated list of ciphers that are not permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information",
			},
			"tls_protocol_allow_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If specified, contains a space-separated list of protocols that are permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information.",
			},
			"tls_protocol_deny_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If specified, contains a space-separated list of protocols that are not permitted for use by TLS clients. See Selecting Ciphers and Protocol for more information.",
			},

			"tls_truststore_filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The filename for the TLS TrustStore used to hold trusted certificates for TLS connections. This can be in the format classpath:path/to/file.p12 or file:///path/to/file.p12. Valid file extensions are .jks (Java Keystore) or .p12 (PKCS#12 store).",
			},
			"tls_truststore_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for the TLS TrustStore (leave blank if the store has no password).",
			},
			// Dependency Options ------------------------
			"dependency_local_inbound_security": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SECURITY_IN_UP",
				Description: "The inbound security module to use for authenticating and authorizing users to this module where authentication requires a username and password.",
			},
			"dependency_fhir_persistence_module": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "PERSISTANCE_ALL",
				Description: "The FHIR Storage module to associate with this module.",
			},
			"dependency_saml_authentication_module": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The SAML Inbound Security module to use when performing a SAML user authentication.",
			},
			"dependency_self_registration_provider_module": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This can be supplied to some interactive modules in order to support self-registration of users.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func smartOutboundSecurityResourceToModuleConfig(d *schema.ResourceData) (*smilecdr.ModuleConfig, error) {

	d.Set("module_type", "SECURITY_OUT_SMART") // Hardcoded for this module type

	moduleConfig := &smilecdr.ModuleConfig{
		ModuleId:   d.Get("module_id").(string),
		ModuleType: d.Get("module_type").(string),
		NodeId:     d.Get("node_id").(string),
	}

	// User Authentication Options ------------------------
	if v, ok := d.GetOk("anonymous_account_username"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "anonymous.access.account_username",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("anonymous_access_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "anonymous.access.enabled",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("saml_authentication_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "saml.enabled",
			Value: v.(string),
		})
	}
	// CORS Options ------------------------
	if v, ok := d.GetOk("cors_allowed_headers"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cors.allowed_headers",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("cors_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cors.enable",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("cors_origins"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cors.origins",
			Value: v.(string),
		})
	}
	// Da Vinci Options ------------------------
	if v, ok := d.GetOk("davinci_native_consent_handling"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "davinci.consent_handling",
			Value: v.(string),
		})
	}
	// HTTP Access Log Options ------------------------
	if v, ok := d.GetOk("http_access_log_appenders"); ok {
		var theValue string = ""
		for _, appender := range v.([]interface{}) {
			theValue = theValue + appender.(string) + "\n"
		}
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "access_log.appenders",
			Value: theValue,
		})
	}
	// HTTP Listener Options ------------------------
	if v, ok := d.GetOk("http_listener_bind_address"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "bind_address",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_listener_context_path"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "context_path",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_listener_endpoint_health_path"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "endpoint_health.path",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_listener_unhealthy_response_code"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "endpoint_health.status_code_if_unhealthy",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_listener_https_forwarding_assumed"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "https_forwarding_assumed",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_listener_port"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "port",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_listener_respect_forward_headers"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "respect_forward_headers",
			Value: v.(string),
		})
	}
	// HTTP Request Pool Options ------------------------
	if v, ok := d.GetOk("http_request_maximum_request_header_size"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "max_header_size.request.kb",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_request_maximum_response_header_size"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "max_header_size.response.kb",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_request_read_idle_timeout"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "read_idle_timeout.millis",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("http_request_thread_pool_accept_queue_size"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "thread_pool.accept_queue_size",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_request_thread_pool_max_size"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "threadpool.max",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_request_thread_pool_min_size"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "threadpool.min",
			Value: v.(string),
		})
	}
	// HTTP Security Options --------------------------------
	if v, ok := d.GetOk("http_security_block_http_head"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "block_http_head",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_security_block_http_options"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "block_http_options",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_security_custom_response_headers"); ok {
		var theValue string = ""
		for _, header := range v.([]interface{}) {
			theValue = theValue + header.(string) + "\n"
		}
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "custom_response_headers",
			Value: theValue,
		})
	}
	if v, ok := d.GetOk("http_security_frame_options_allow_from"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "frame_options.allow_from",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_security_pin_host"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "pin_host",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_security_suppress_error_details"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "suppress_error_details",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("http_security_suppress_platform_info"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "suppress_platform_info",
			Value: v.(string),
		})
	}
	// JavaScript Execution Environment Options ------------------------
	if v, ok := d.GetOk("javascript_debug_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.debug_enabled",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("javascript_debug_host_address"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.host_address",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("javascript_debug_path"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.path",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("javascript_debug_port"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.port",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("javascript_debug_secure"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.secure",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("javascript_debug_suspend"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.suspend",
			Value: v.(string),
		})
	}
	// JWKS Options --------------------------------
	if v, ok := d.GetOk("jwks_keystore_id"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "openid.signing.keystore_id",
			Value: v.(string),
		})
	}
	// OIDC Token Validation Options ------------------------
	if v, ok := d.GetOk("oidc_http_client_jwks_cache_timeout"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_client.jwks_cache.mins",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("oidc_http_client_truststore_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_client.truststore.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("oidc_http_client_truststore_password"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_client.truststore.password",
			Value: v.(string),
		})
	}
	// OpenID Connect (OIDC) Options ------------------------
	if v, ok := d.GetOk("oidc_pkce_required"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "pkce.required",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("oidc_pkce_plain_challenge_supported"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "pkce.plain_challenge_supported",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("oidc_cache_authorization_tokens"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cache.authorized_tokens.millis",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("oidc_client_secret_encoding"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "client_secret.encoding",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("oidc_client_secret_expiry_duration"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "client_secret.expiry_duration_days",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("oidc_issuer_url"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "issuer.url",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("oidc_rotate_token_after_use"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "rotate_refresh_token_after_use",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("oidc_smart_capabilities_list"); ok {
		var theValue string = ""
		for _, capability := range v.([]interface{}) {
			theValue = theValue + capability.(string) + "\n"
		}
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "smart_capabilities_list",
			Value: theValue,
		})
	}
	// OAuth2/OIDC Federation Options ------------------------
	if v, ok := d.GetOk("oauth2_federation_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "federate_mode.enabled",
			Value: v.(string),
		})
	}
	// SMART Callback Script Options ------------------------
	if v, ok := d.GetOk("smart_callback_post_authorize_script_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "post_authorize_script.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_callback_post_authorize_script_text"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "post_authorize_script.text",
			Value: v.(string),
		})
	}
	// CODAP Options ------------------------
	if v, ok := d.GetOk("codap_authorization_script_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "codap.auth_script.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("codap_authorization_script_text"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "codap.auth_script.text",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("codap_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "codap.enabled",
			Value: v.(string),
		})
	}
	// SMART Login Skin Options ------------------------
	if v, ok := d.GetOk("smart_login_skin_approval_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.approve_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_context_selection_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.context_selection.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_error_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.error_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_user_registration_forgot_password_template_step1"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.forgot_password_step1.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_user_registration_forgot_password_template_step2"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.forgot_password_step2.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_user_registration_forgot_password_template_step3"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.forgot_password_step3.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_federated_oath2_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.login_oauth2_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_login_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.login_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_user_registration_template_step1"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.register_step1.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_user_registration_template_step2"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.register_step2.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_session_management_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.session_management_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_2fa_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.tfa_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_terms_of_service_template"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.tos_page.template",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_login_skin_webjar_id"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "skin.webjar_id",
			Value: v.(string),
		})
	}
	// SMART Login Terms of Service Options ------------------------
	if v, ok := d.GetOk("smart_login_terms_of_service"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tos.version_string",
			Value: v.(string),
		})
	}
	// SMART Authorization Options ------------------------
	if v, ok := d.GetOk("smart_authorization_allowed_audience_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "allowed_audience_list",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_authorization_email_from_address"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "email.from_address",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_authorization_enforce_approved_scopes"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "enforce_approved_scopes_to_restrict_permissions",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_authorization_scopes_supported"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "smart_configuration.scopes_supported",
			Value: v.(string),
		})
	}
	// SMART Definitions Seeding Options ------------------------
	if v, ok := d.GetOk("openid_connect_server_pre_seed_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "seed_clients.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("openid_connect_client_pre_seed_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "seed_servers.file",
			Value: v.(string),
		})
	}
	// Sessions Options ------------------------
	if v, ok := d.GetOk("sessions_in_memory"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "sessions.inmemory",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("sessions_max_concurrent_sessions_per_user"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "sessions.maximum_concurrent",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("sessions_scavenger_interval_ms"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "sessions.scavenger.interval.millis",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("sessions_timeout_mins"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "sessions.timeout.mins",
			Value: v.(string),
		})
	}
	// TLS Options --------------------------------
	if v, ok := d.GetOk("tls_client_auth_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.clientauth.enabled",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.enabled",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_keystore_filename"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.keystore.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_keystore_key_alias"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.keystore.keyalias",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_keystore_key_password"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.keystore.keypass",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_keystore_password"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.keystore.password",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_cipher_allow_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.protocol.cipher_whitelist",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_cipher_deny_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.protocol.cipher_blacklist",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_protocol_allow_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.protocol.protocol_whitelist",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_protocol_deny_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.protocol.protocol_blacklist",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_truststore_filename"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.truststore.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tls_truststore_password"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tls.truststore.password",
			Value: v.(string),
		})
	}
	// Dependencies --------------------------------
	if v, ok := d.GetOk("dependency_local_inbound_security"); ok {
		moduleConfig.Dependencies = append(moduleConfig.Dependencies, smilecdr.ModuleDependencies{
			ModuleId: "local_security",
			Type:     v.(string),
		})

	}
	if v, ok := d.GetOk("dependency_fhir_persistence_module"); ok {
		moduleConfig.Dependencies = append(moduleConfig.Dependencies, smilecdr.ModuleDependencies{
			ModuleId: "persistence",
			Type:     v.(string),
		})
	}
	if v, ok := d.GetOk("dependency_saml_authentication_module"); ok {
		moduleConfig.Dependencies = append(moduleConfig.Dependencies, smilecdr.ModuleDependencies{
			ModuleId: "security_in_saml",
			Type:     v.(string),
		})
	}
	if v, ok := d.GetOk("dependency_self_service_user_management_module"); ok {
		moduleConfig.Dependencies = append(moduleConfig.Dependencies, smilecdr.ModuleDependencies{
			ModuleId: "self_service_user_management",
			Type:     v.(string),
		})
	}
	return moduleConfig, nil
}

func resourceSmartOutboundSecurityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleConfig, err := smartOutboundSecurityResourceToModuleConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}

	module, err := c.PostModuleConfig(*moduleConfig)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(module.ModuleId) // the primary resource identifier. must be unique.

	return resourceSmartOutboundSecurityRead(ctx, d, m)

}

func resourceSmartOutboundSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)
	moduleConfig, err := c.GetModuleConfig(nodeId, moduleId)

	// map from moduleConfig to resourceData
	if err != nil {
		return diag.FromErr(err)
	}

	// User Authentication Options ------------------------
	val, ok := moduleConfig.LookupOptionOk("anonymous.access.account_username")
	if ok {
		d.Set("anonymous_account_username", val)
	}
	val, ok = moduleConfig.LookupOptionOk("anonymous.access.enabled")
	if ok {
		d.Set("anonymous_access_enabled", val)
	}
	val, ok = moduleConfig.LookupOptionOk("saml.enabled")
	if ok {
		d.Set("saml_authentication_enabled", val)
	}
	val, ok = moduleConfig.LookupOptionOk("cors.allowed_headers")
	if ok {
		d.Set("cors_allowed_headers", val)
	}
	val, ok = moduleConfig.LookupOptionOk("cors.enable")
	if ok {
		d.Set("cors_enabled", val)
	}
	val, ok = moduleConfig.LookupOptionOk("cors.origins")
	if ok {
		d.Set("cors_origins", val)
	}
	// Da Vinci Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("davinci.consent_handling")
	if ok {
		d.Set("davinci_native_consent_handling", val)
	}
	// HTTP Access Log Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("access_log.appenders")
	if ok {
		appenders := strings.Split(val, "\n")
		d.Set("http_access_log_appenders", appenders)
	}
	// HTTP Listener Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("bind_address")
	if ok {
		d.Set("http_listener_bind_address", val)
	}
	val, ok = moduleConfig.LookupOptionOk("context_path")
	if ok {
		d.Set("http_listener_context_path", val)
	}
	val, ok = moduleConfig.LookupOptionOk("endpoint_health.path")
	if ok {
		d.Set("http_listener_endpoint_health_path", val)
	}
	val, ok = moduleConfig.LookupOptionOk("endpoint_health.status_code_if_unhealthy")
	if ok {
		d.Set("http_listener_unhealthy_response_code", val)
	}
	val, ok = moduleConfig.LookupOptionOk("https_forwarding_assumed")
	if ok {
		d.Set("http_listener_https_forwarding_assumed", val)
	}
	val, ok = moduleConfig.LookupOptionOk("port")
	if ok {
		d.Set("http_listener_port", val)
	}
	val, ok = moduleConfig.LookupOptionOk("respect_forward_headers")
	if ok {
		d.Set("http_listener_respect_forward_headers", val)
	}
	// HTTP Request Pool Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("max_header_size.request.kb")
	if ok {
		d.Set("http_request_maximum_request_header_size", val)
	}
	val, ok = moduleConfig.LookupOptionOk("max_header_size.response.kb")
	if ok {
		d.Set("http_request_maximum_response_header_size", val)
	}
	val, ok = moduleConfig.LookupOptionOk("read_idle_timeout.millis")
	if ok {
		d.Set("http_request_read_idle_timeout", val)
	}
	val, ok = moduleConfig.LookupOptionOk("thread_pool.accept_queue_size")
	if ok {
		d.Set("http_request_thread_pool_accept_queue_size", val)
	}
	val, ok = moduleConfig.LookupOptionOk("threadpool.max")
	if ok {
		d.Set("http_request_thread_pool_max_size", val)
	}
	val, ok = moduleConfig.LookupOptionOk("threadpool.min")
	if ok {
		d.Set("http_request_thread_pool_min_size", val)
	}
	// HTTP Security Options --------------------------------
	val, ok = moduleConfig.LookupOptionOk("block_http_head")
	if ok {
		d.Set("http_security_block_http_head", val)
	}
	val, ok = moduleConfig.LookupOptionOk("custom_response_headers")
	if ok {
		headers := strings.Split(val, "\n")
		d.Set("http_security_custom_response_headers", headers)
	}
	val, ok = moduleConfig.LookupOptionOk("frame_options.allow_from")
	if ok {
		d.Set("http_security_frame_options_allow_from", val)
	}
	val, ok = moduleConfig.LookupOptionOk("pin_host")
	if ok {
		d.Set("http_security_pin_host", val)
	}
	val, ok = moduleConfig.LookupOptionOk("suppress_error_details")
	if ok {
		d.Set("http_security_suppress_error_details", val)
	}
	val, ok = moduleConfig.LookupOptionOk("suppress_platform_info")
	if ok {
		d.Set("http_security_suppress_platform_info", val)
	}
	// JavaScript Execution Environment Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("debug.debug_enabled")
	if ok {
		d.Set("javascript_debug_enabled", val)
	}
	val, ok = moduleConfig.LookupOptionOk("debug.host_address")
	if ok {
		d.Set("javascript_debug_host_address", val)
	}
	val, ok = moduleConfig.LookupOptionOk("debug.path")
	if ok {
		d.Set("javascript_debug_path", val)
	}
	val, ok = moduleConfig.LookupOptionOk("debug.port")
	if ok {
		d.Set("javascript_debug_port", val)
	}
	val, ok = moduleConfig.LookupOptionOk("debug.secure")
	if ok {
		d.Set("javascript_debug_secure", val)
	}
	val, ok = moduleConfig.LookupOptionOk("debug.suspend")
	if ok {
		d.Set("javascript_debug_suspend", val)
	}
	// JWKS Options --------------------------------
	val, ok = moduleConfig.LookupOptionOk("openid.signing.keystore_id")
	if ok {
		d.Set("jwks_keystore_id", val)
	}
	// OIDC Token Validation Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("introspection_client.jwks_cache.mins")
	if ok {
		d.Set("oidc_http_client_jwks_cache_timeout", val)
	}
	val, ok = moduleConfig.LookupOptionOk("introspection_client.truststore.file")
	if ok {
		d.Set("oidc_http_client_truststore_file", val)
	}
	val, ok = moduleConfig.LookupOptionOk("introspection_client.truststore.password")
	if ok {
		d.Set("oidc_http_client_truststore_password", val)
	}
	// OpenID Connect (OIDC) Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("pkce.plain_challenge_supported")
	if ok {
		d.Set("oidc_pkce_plain_challenge_supported", val)
	}
	val, ok = moduleConfig.LookupOptionOk("pkce.required")
	if ok {
		d.Set("oidc_pkce_required", val)
	}
	val, ok = moduleConfig.LookupOptionOk("cache.authorized_tokens.millis")
	if ok {
		d.Set("oidc_cache_authorization_tokens", val)
	}
	val, ok = moduleConfig.LookupOptionOk("client_secret.encoding")
	if ok {
		d.Set("oidc_client_secret_encoding", val)
	}
	val, ok = moduleConfig.LookupOptionOk("client_secret.expiry_duration_days")
	if ok {
		d.Set("oidc_client_secret_expiry_duration", val)
	}
	val, ok = moduleConfig.LookupOptionOk("issuer.url")
	if ok {
		d.Set("oidc_issuer_url", val)
	}
	val, ok = moduleConfig.LookupOptionOk("rotate_refresh_token_after_use")
	if ok {
		d.Set("oidc_rotate_token_after_use", val)
	}
	val, ok = moduleConfig.LookupOptionOk("smart_capabilities_list")
	if ok {
		capabilities := strings.Split(val, "\n")
		d.Set("oidc_smart_capabilities_list", capabilities)
	}
	// OAuth2/OIDC Federation Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("federate_mode.enabled")
	if ok {
		d.Set("oauth2_federation_enabled", val)
	}
	// SMART Callback Script Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("post_authorize_script.file")
	if ok {
		d.Set("smart_callback_post_authorize_script_file", val)
	}
	val, ok = moduleConfig.LookupOptionOk("post_authorize_script.text")
	if ok {
		d.Set("smart_callback_post_authorize_script_text", val)
	}
	// CODAP Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("codap.auth_script.file")
	if ok {
		d.Set("codap_authorization_script_file", val)
	}
	val, ok = moduleConfig.LookupOptionOk("codap.auth_script.text")
	if ok {
		d.Set("codap_authorization_script_text", val)
	}
	val, ok = moduleConfig.LookupOptionOk("codap.enabled")
	if ok {
		d.Set("codap_enabled", val)
	}
	// SMART Login Skin Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("skin.approve_page.template")
	if ok {
		d.Set("smart_login_skin_approval_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.context_selection.template")
	if ok {
		d.Set("smart_login_skin_context_selection_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.error_page.template")
	if ok {
		d.Set("smart_login_skin_error_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.forgot_password_step1.template")
	if ok {
		d.Set("smart_login_skin_user_registration_forgot_password_template_step1", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.forgot_password_step2.template")
	if ok {
		d.Set("smart_login_skin_user_registration_forgot_password_template_step2", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.forgot_password_step3.template")
	if ok {
		d.Set("smart_login_skin_user_registration_forgot_password_template_step3", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.login_oauth2_page.template")
	if ok {
		d.Set("smart_login_skin_federated_oath2_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.login_page.template")
	if ok {
		d.Set("smart_login_skin_login_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.register_step1.template")
	if ok {
		d.Set("smart_login_skin_user_registration_template_step1", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.register_step2.template")
	if ok {
		d.Set("smart_login_skin_user_registration_template_step2", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.session_management_page.template")
	if ok {
		d.Set("smart_login_skin_session_management_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.tfa_page.template")
	if ok {
		d.Set("smart_login_skin_2fa_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.tos_page.template")
	if ok {
		d.Set("smart_login_skin_terms_of_service_template", val)
	}
	val, ok = moduleConfig.LookupOptionOk("skin.webjar_id")
	if ok {
		d.Set("smart_login_skin_webjar_id", val)
	}
	// SMART Login Terms of Service Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("tos.version_string")
	if ok {
		d.Set("smart_login_terms_of_service", val)
	}
	// SMART Authorization Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("allowed_audience_list")
	if ok {
		d.Set("smart_authorization_allowed_audience_list", val)
	}
	val, ok = moduleConfig.LookupOptionOk("email.from_address")
	if ok {
		d.Set("smart_authorization_email_from_address", val)
	}
	val, ok = moduleConfig.LookupOptionOk("enforce_approved_scopes_to_restrict_permissions")
	if ok {
		d.Set("smart_authorization_enforce_approved_scopes", val)
	}
	val, ok = moduleConfig.LookupOptionOk("smart_configuration.scopes_supported")
	if ok {
		d.Set("smart_authorization_scopes_supported", val)
	}
	// SMART Definitions Seeding Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("seed_clients.file")
	if ok {
		d.Set("openid_connect_server_pre_seed_file", val)
	}
	val, ok = moduleConfig.LookupOptionOk("seed_servers.file")
	if ok {
		d.Set("openid_connect_client_pre_seed_file", val)
	}
	// Sessions Options ------------------------
	val, ok = moduleConfig.LookupOptionOk("sessions.inmemory")
	if ok {
		d.Set("sessions_in_memory", val)
	}
	val, ok = moduleConfig.LookupOptionOk("sessions.maximum_concurrent")
	if ok {
		d.Set("sessions_max_concurrent_sessions_per_user", val)
	}
	val, ok = moduleConfig.LookupOptionOk("sessions.scavenger.interval.millis")
	if ok {
		d.Set("sessions_scavenger_interval_ms", val)
	}
	val, ok = moduleConfig.LookupOptionOk("sessions.timeout.mins")
	if ok {
		d.Set("sessions_timeout_mins", val)
	}
	// TLS Options --------------------------------
	val, ok = moduleConfig.LookupOptionOk("tls.clientauth.enabled")
	if ok {
		d.Set("tls_client_auth_enabled", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.enabled")
	if ok {
		d.Set("tls_enabled", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.keystore.file")
	if ok {
		d.Set("tls_keystore_filename", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.keystore.keyalias")
	if ok {
		d.Set("tls_keystore_key_alias", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.keystore.keypass")
	if ok {
		d.Set("tls_keystore_key_password", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.keystore.password")
	if ok {
		d.Set("tls_keystore_password", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.protocol.cipher_whitelist")
	if ok {
		d.Set("tls_cipher_allow_list", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.protocol.cipher_blacklist")
	if ok {
		d.Set("tls_cipher_deny_list", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.protocol.protocol_whitelist")
	if ok {
		d.Set("tls_protocol_allow_list", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.protocol.protocol_blacklist")
	if ok {
		d.Set("tls_protocol_deny_list", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.truststore.file")
	if ok {
		d.Set("tls_truststore_filename", val)
	}
	val, ok = moduleConfig.LookupOptionOk("tls.truststore.password")
	if ok {
		d.Set("tls_truststore_password", val)
	}

	// Set The Specific Dependencies for SMART Outbound Security
	for _, dependency := range moduleConfig.Dependencies {
		if dependency.ModuleId == "local_security" {
			d.Set("dependency_local_inbound_security", dependency.Type)
		}
		if dependency.ModuleId == "persistence" {
			d.Set("dependency_fhir_persistence_module", dependency.Type)
		}
		if dependency.ModuleId == "security_in_saml" {
			d.Set("dependency_saml_authentication_module", dependency.Type)
		}
		if dependency.ModuleId == "self_service_user_management" {
			d.Set("dependency_self_service_user_management_module", dependency.Type)
		}
	}

	return nil
}

func resourceSmartOutboundSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleConfig, err := smartOutboundSecurityResourceToModuleConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(moduleConfig.ModuleId) // the primary resource identifier. must be unique.

	_, pErr := c.PutModuleConfig(*moduleConfig)

	if pErr != nil {
		return diag.FromErr(pErr)
	}

	return resourceSmartOutboundSecurityRead(ctx, d, m)
}

func resourceSmartOutboundSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*smilecdr.Client)

	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)

	err := c.DeleteModuleConfig(nodeId, moduleId)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("") // This is unset when the resource is deleted

	return nil
}
