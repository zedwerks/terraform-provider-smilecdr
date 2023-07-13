// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zed-werks/terraform-smilecdr/provider/util"
	"github.com/zed-werks/terraform-smilecdr/smilecdr"
)

var (
	smileCdrUserPermissionTypes = []string{"ACCESS_ADMIN_JSON", "ACCESS_ADMIN_WEB",
		"ACCESS_FHIRWEB",
		"ACCESS_FHIR_ENDPOINT",
		"AG_ADMIN_CONSOLE_READ",
		"AG_ADMIN_CONSOLE_WRITE",
		"AG_DEV_PORTAL_READ",
		"AG_DEV_PORTAL_WRITE",
		"ARCHIVE_MODULE",
		"BLOCK_FHIR_READ_UNLESS_CODE_IN_VS",
		"BLOCK_FHIR_READ_UNLESS_CODE_NOT_IN_VS",
		"CDA_IMPORT",
		"CHANGE_OWN_DEFAULT_LAUNCH_CONTEXTS",
		"CHANGE_OWN_PASSWORD",
		"CHANGE_OWN_TFA_KEY",
		"CONTROL_MODULE",
		"CREATE_CDA_TEMPLATE",
		"CREATE_MODULE",
		"CREATE_USER",
		"DELETE_CDA_TEMPLATE",
		"DOCREF",
		"DQM_QPP_BUILD",
		"EMPI_ADMIN",
		"EMPI_UPDATE_MATCH_RULES",
		"EMPI_VIEW_MATCH_RULES",
		"ETL_IMPORT_PROCESS_FILE",
		"FHIR_ACCESS_PARTITION_ALL",
		"FHIR_ACCESS_PARTITION_NAME",
		"FHIR_ALL_DELETE",
		"FHIR_ALL_READ",
		"FHIR_ALL_WRITE",
		"FHIR_AUTO_MDM",
		"FHIR_BATCH",
		"FHIR_CAPABILITIES",
		"FHIR_DELETE_ALL_IN_COMPARTMENT",
		"FHIR_DELETE_ALL_OF_TYPE",
		"FHIR_DELETE_CASCADE_ALLOWED",
		"FHIR_DELETE_EXPUNGE",
		"FHIR_DELETE_TYPE_IN_COMPARTMENT",
		"FHIR_DTR_USER",
		"FHIR_EMPI_ADMIN",
		"FHIR_EXPUNGE_DELETED",
		"FHIR_EXPUNGE_EVERYTHING",
		"FHIR_EXPUNGE_PREVIOUS_VERSIONS",
		"FHIR_EXTENDED_OPERATION_ON_ANY_INSTANCE",
		"FHIR_EXTENDED_OPERATION_ON_ANY_INSTANCE_OF_TYPE",
		"FHIR_EXTENDED_OPERATION_ON_SERVER",
		"FHIR_EXTENDED_OPERATION_ON_TYPE",
		"FHIR_GET_RESOURCE_COUNTS",
		"FHIR_GRAPHQL",
		"FHIR_LIVEBUNDLE",
		"FHIR_MANAGE_PARTITIONS",
		"FHIR_MANUAL_VALIDATION",
		"FHIR_MDM_ADMIN",
		"FHIR_META_OPERATIONS_SUPERUSER",
		"FHIR_MODIFY_SEARCH_PARAMETERS",
		"FHIR_OP_APPLY",
		"FHIR_OP_BINARY_ACCESS_READ",
		"FHIR_OP_BINARY_ACCESS_WRITE",
		"FHIR_OP_CQL_EVALUATE_MEASURE",
		"FHIR_OP_EMPI_CLEAR",
		"FHIR_OP_EMPI_DUPLICATE_PERSONS",
		"FHIR_OP_EMPI_MERGE_PERSONS",
		"FHIR_OP_EMPI_QUERY_LINKS",
		"FHIR_OP_EMPI_SUBMIT",
		"FHIR_OP_EMPI_UPDATE_LINK",
		"FHIR_OP_ENCOUNTER_EVERYTHING",
		"FHIR_OP_EXTRACT",
		"FHIR_OP_INITIATE_BULK_DATA_EXPORT",
		"FHIR_OP_INITIATE_BULK_DATA_EXPORT_GROUP",
		"FHIR_OP_INITIATE_BULK_DATA_EXPORT_PATIENT",
		"FHIR_OP_INITIATE_BULK_DATA_EXPORT_SYSTEM",
		"FHIR_OP_INITIATE_BULK_DATA_IMPORT",
		"FHIR_OP_MDM_CLEAR",
		"FHIR_OP_MDM_CREATE_LINK",
		"FHIR_OP_MDM_DUPLICATE_GOLDEN_RESOURCES",
		"FHIR_OP_MDM_LINK_HISTORY",
		"FHIR_OP_MDM_MERGE_GOLDEN_RESOURCES",
		"FHIR_OP_MDM_NOT_DUPLICATE",
		"FHIR_OP_MDM_QUERY_LINKS",
		"FHIR_OP_MDM_SUBMIT",
		"FHIR_OP_MDM_UPDATE_LINK",
		"FHIR_OP_MEMBER_MATCH",
		"FHIR_OP_PACKAGE",
		"FHIR_OP_PATIENT_EVERYTHING",
		"FHIR_OP_PATIENT_MATCH",
		"FHIR_OP_PATIENT_SUMMARY",
		"FHIR_OP_POPULATE",
		"FHIR_OP_PREPOPULATE",
		"FHIR_OP_STRUCTUREDEFINITION_SNAPSHOT",
		"FHIR_PATCH",
		"FHIR_PROCESS_MESSAGE",
		"FHIR_READ_ALL_IN_COMPARTMENT",
		"FHIR_READ_ALL_OF_TYPE",
		"FHIR_READ_INSTANCE",
		"FHIR_READ_SEARCH_PARAMETERS",
		"FHIR_READ_TYPE_IN_COMPARTMENT",
		"FHIR_TRANSACTION",
		"FHIR_TRIGGER_SUBSCRIPTION",
		"FHIR_UPDATE_REWRITE_HISTORY",
		"FHIR_UPLOAD_EXTERNAL_TERMINOLOGY",
		"FHIR_WRITE_ALL_IN_COMPARTMENT",
		"FHIR_WRITE_ALL_OF_TYPE",
		"FHIR_WRITE_INSTANCE",
		"FHIR_WRITE_TYPE_IN_COMPARTMENT",
		"INVOKE_CDS_HOOKS",
		"MANAGE_BATCH_JOBS",
		"MDM_ADMIN",
		"MDM_UPDATE_MATCH_RULES",
		"MDM_VIEW_MATCH_RULES",
		"MODULE_ADMIN",
		"OIDC_CLIENT_PRESET_PERMISSION",
		"OPENID_CONNECT_ADD_CLIENT",
		"OPENID_CONNECT_ADD_SERVER",
		"OPENID_CONNECT_EDIT_CLIENT",
		"OPENID_CONNECT_EDIT_SERVER",
		"OPENID_CONNECT_MANAGE_GLOBAL_SESSIONS",
		"OPENID_CONNECT_VIEW_CLIENT_LIST",
		"OPENID_CONNECT_VIEW_SERVER_LIST",
		"PACKAGE_REGISTRY_READ",
		"PACKAGE_REGISTRY_WRITE",
		"REINSTATE_MODULE",
		"ROLE_ANONYMOUS",
		"ROLE_FHIR_CLIENT",
		"ROLE_FHIR_CLIENT_SUPERUSER",
		"ROLE_FHIR_CLIENT_SUPERUSER_RO",
		"ROLE_FHIR_TERMINOLOGY_READ_CLIENT",
		"ROLE_SUPERUSER",
		"ROLE_SYSTEM",
		"ROLE_SYSTEM_INITIALIZATION",
		"SAVE_USER",
		"START_STOP_MODULE",
		"SUBMIT_ATTACHMENT",
		"UPDATE_MODULE_CONFIG",
		"UPDATE_USER",
		"USE_CDA_TEMPLATE",
		"VIEW_AUDIT_LOG",
		"VIEW_BATCH_JOBS",
		"VIEW_CDA_TEMPLATE",
		"VIEW_METRICS",
		"VIEW_MODULE_CONFIG",
		"VIEW_MODULE_STATUS",
		"VIEW_TRANSACTION_LOG",
		"VIEW_TRANSACTION_LOG_EVENT",
		"VIEW_USERS"}
	smileCdrOpenIdAuthorizationFlows = []string{"AUTHORIZATION_CODE", "CLIENT_CREDENTIALS", "IMPLICIT", "JWT_BEARER", "PASSWORD", "REFRESH_TOKEN"}
)

func resourceOpenIdClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenIdClientCreate,
		ReadContext:   resourceOpenIdClientRead,
		UpdateContext: resourceOpenIdClientUpdate,
		DeleteContext: resourceOpenIdClientDelete,
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "Master",
			},
			"module_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "smart_auth",
			},
			"access_token_validity_seconds": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Default:  300,
			},
			"allowed_grant_types": {
				Type:     schema.TypeSet,
				Required: false,
				Optional: true,
				Default:  []string{"AUTHORIZATION_CODE", "REFRESH_TOKEN"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Required:     false,
					ValidateFunc: schema.SchemaValidateFunc(validation.StringInSlice(smileCdrOpenIdAuthorizationFlows, false)),
				},
			},
			"auto_approve_scopes": {
				Type:     schema.TypeSet,
				Required: false,
				Optional: true,
				Default:  []string{"openid"},
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
				},
			},
			"auto_grant_scopes": {
				Type:     schema.TypeSet,
				Required: false,
				Optional: true,
				Default:  []string{"openid"},
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
				},
			},
			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: util.Validateclient_id,
			},
			"client_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"clientSecrets": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret": {
							Type:     schema.TypeString,
							Required: false,
							Default:  "",
						},
						"description": {
							Type:     schema.TypeString,
							Required: false,
							Default:  "",
						},
						"activation": {
							Type:         schema.TypeString,
							Required:     false,
							Default:      "",
							ValidateFunc: validation.IsRFC3339Time,
						},
						"expiration": {
							Type:         schema.TypeString,
							Required:     false,
							Default:      "",
							ValidateFunc: validation.IsRFC3339Time,
						},
					},
				},
			},
			"fixed_scope": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"refresh_token_validity_seconds": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Default:  86400,
			},
			"registered_redirect_uris": {
				Type:     schema.TypeSet,
				Required: false,
				Optional: true,
				Default:  []string{},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scopes": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Default:  []string{"openid"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secret_required": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"secret_client_can_change": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"can_introspect_any_tokens": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"can_introspect_own_tokens": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"always_require_approval": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"can_reissue_tokens": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Required: false,
				Optional: true,
				Default:  []interface{}{},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: schema.SchemaValidateFunc(validation.StringInSlice(smileCdrUserPermissionTypes, false)),
						},
						"argument": {
							Type:     schema.TypeString,
							Required: false,
						},
					},
				},
			},
			"remember_approved_scopes": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"attestation_accepted": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"public_jwks": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
			},
			"jwks_url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
			},
			"archivedAt": {
				Type:         schema.TypeString,
				Required:     false,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"created_by_app_sphere": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDataToOpenIdClient(d *schema.ResourceData) *smilecdr.OpenIdClient {

	secrets := d.Get("client_secrets").(*schema.Set).List()
	clientSecrets := make([]smilecdr.ClientSecret, len(secrets))

	permissions := d.Get("permissions").(*schema.Set).List()
	userPermissions := make([]smilecdr.UserPermission, len(permissions))

	openidClient := &smilecdr.OpenIdClient{
		ClientId:                    d.Get("client_id").(string),
		ClientName:                  d.Get("client_name").(string),
		NodeId:                      d.Get("node_id").(string),
		ModuleId:                    d.Get("module_id").(string),
		AccessTokenValiditySeconds:  d.Get("access_token_validity_seconds").(int),
		AllowedGrantTypes:           d.Get("allowedGrantTypes").([]string),
		AutoApproveScopes:           d.Get("auto_approve_scopes").([]string),
		AutoGrantScopes:             d.Get("auto_grant_scopes").([]string),
		ClientSecrets:               clientSecrets,
		FixedScope:                  d.Get("fixed_scope").(bool),
		RefreshTokenValiditySeconds: d.Get("refresh_token_validity_seconds").(int),
		RegisteredRedirectUris:      d.Get("registered_redirect_uris").([]string),
		Scopes:                      d.Get("scopes").([]string),
		SecretRequired:              d.Get("secret_required").(bool),
		SecretClientCanChange:       d.Get("secret_client_can_change").(bool),
		Enabled:                     d.Get("enabled").(bool),
		CanIntrospectAnyTokens:      d.Get("can_introspect_any_tokens").(bool),
		CanIntrospectOwnTokens:      d.Get("can_introspect_own_tokens").(bool),
		AlwaysRequireApproval:       d.Get("always_require_approval").(bool),
		CanReissueTokens:            d.Get("can_reissue_tokens").(bool),
		Permissions:                 userPermissions,
		AttestationAccepted:         d.Get("attestation_accepted").(bool),
		PublicJwksUri:               d.Get("public_jwksUri").(string),
		ArchivedAt:                  d.Get("archivedAt").(string),
	}
	return openidClient

}

func resourceOpenIdClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	client := resourceDataToOpenIdClient(d)

	_, err := c.PostOpenIdClient(*client)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags

}

func resourceOpenIdClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	client_id := d.Get("client_id").(string)
	nodeId := d.Get("node_id").(string)
	moduleId := d.Get("module_id").(string)

	openIdClient, err := c.GetOpenIdClient(client_id, nodeId, moduleId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(openIdClient.ClientId)
	d.Set("client_name", openIdClient.ClientName)
	d.Set("node_id", openIdClient.NodeId)
	d.Set("module_id", openIdClient.ModuleId)
	d.Set("access_token_validity_seconds", openIdClient.AccessTokenValiditySeconds)
	d.Set("allowed_grant_types", openIdClient.AllowedGrantTypes)
	d.Set("auto_approve_scopes", openIdClient.AutoApproveScopes)
	d.Set("auto_grant_scopes", openIdClient.AutoGrantScopes)
	d.Set("client_secrets", openIdClient.ClientSecrets)
	d.Set("fixed_scope", openIdClient.FixedScope)
	d.Set("refresh_token_validity_seconds", openIdClient.RefreshTokenValiditySeconds)
	d.Set("registered_redirect_uris", openIdClient.RegisteredRedirectUris)
	d.Set("scopes", openIdClient.Scopes)
	d.Set("secret_required", openIdClient.SecretRequired)
	d.Set("secret_client_can_change", openIdClient.SecretClientCanChange)
	d.Set("enabled", openIdClient.Enabled)
	d.Set("can_introspect_any_tokens", openIdClient.CanIntrospectAnyTokens)
	d.Set("can_introspect_own_tokens", openIdClient.CanIntrospectOwnTokens)
	d.Set("always_require_approval", openIdClient.AlwaysRequireApproval)
	d.Set("can_reissue_tokens", openIdClient.CanReissueTokens)
	d.Set("permissions", openIdClient.Permissions)
	d.Set("attestation_accepted", openIdClient.AttestationAccepted)
	d.Set("public_jwks_uri", openIdClient.PublicJwksUri)
	d.Set("archived_at", openIdClient.ArchivedAt)
	d.Set("created_by_app_sphere", openIdClient.CreatedByAppSphere)
	return diags

}

func resourceOpenIdClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	client := resourceDataToOpenIdClient(d)

	_, err := c.PutOpenIdClient(*client)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags

}

func resourceOpenIdClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}
