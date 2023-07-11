// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"pid": {
				Type:     schema.TypeInt,
				Required: false,
			},
			"nodeId": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "Master",
			},
			"moduleId": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "smart_auth",
			},
			"accessTokenValiditySeconds": {
				Type:     schema.TypeInt,
				Required: false,
				Default:  300,
			},
			"allowedGrantTypes": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Required:     false,
					ValidateFunc: schema.SchemaValidateFunc(validation.StringInSlice(smileCdrOpenIdAuthorizationFlows, false)),
				},
			},
			"autoApproveScopes": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
				},
			},
			"autoGrantScopes": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
				},
			},
			"clientId": {
				Type:     schema.TypeString,
				Required: true,
			},
			"clientName": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "Some Client",
			},
			"clientSecrets": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"pid": {
							Type:     schema.TypeInt,
							Required: false,
						},
						"secret": {
							Type:     schema.TypeString,
							Required: true,
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
			"fixedScope": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"refreshTokenValiditySeconds": {
				Type:     schema.TypeInt,
				Required: false,
				Default:  86400,
			},
			"registeredRedirectUris": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scopes": {
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"secretRequired": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"secretClientCanChange": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  true,
			},
			"canIntrospectAnyTokens": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"canIntrospectOwnTokens": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"alwaysRequireApproval": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"canReissueTokens": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"permissions": {
				Type:     schema.TypeSet,
				Required: false,
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
			"rememberApprovedScopes": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"attestationAccepted": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
			"publicJwks": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "",
			},
			"jwksUrl": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "",
			},
			"archivedAt": {
				Type:         schema.TypeString,
				Required:     false,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"createdByAppSphere": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
			},
		},
	}
}

func resourceOpenIdClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}

func resourceOpenIdClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}

func resourceOpenIdClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}

func resourceOpenIdClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}
