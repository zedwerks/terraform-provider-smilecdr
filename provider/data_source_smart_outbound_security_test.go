// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestSmartOutboundSecurityDataSource(t *testing.T) {

	t.Parallel()
	moduleName := "smart_" + acctest.RandString(8)
	dataSourceName := "data.smilecdr_smart_outbound_security.testacc"
	resourceName := "smilecdr_smart_outbound_security.testacc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSmartOutboundDataSource_Config(moduleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "module_id", resourceName, "module_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_id", resourceName, "node_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "http_listener_context_path", resourceName, "http_listener_context_path"),
					resource.TestCheckResourceAttrPair(dataSourceName, "http_listener_unhealthy_response_code", resourceName, "http_listener_unhealthy_response_code"),
					resource.TestCheckResourceAttrPair(dataSourceName, "http_listener_port", resourceName, "http_listener_port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "oidc_issuer_url", resourceName, "oidc_issuer_url"),
					resource.TestCheckResourceAttrPair(dataSourceName, "oidc_federate_mode_enabled", resourceName, "oidc_federate_mode_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "smart_authorization_enforce_approved_scopes", resourceName, "smart_authorization_enforce_approved_scopes"),
					resource.TestCheckResourceAttrPair(dataSourceName, "smart_authorization_scopes_supported", resourceName, "smart_authorization_scopes_supported"),
					resource.TestCheckResourceAttrPair(dataSourceName, "sessions_max_concurrent_sessions_per_user", resourceName, "sessions_max_concurrent_sessions_per_user"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dependency_fhir_persistence_module", resourceName, "dependency_fhir_persistence_module"),
				),
			},
		},
	})
}

func testSmartOutboundDataSource_Config(moduleName string) string {

	return fmt.Sprintf(`resource "smilecdr_smart_outbound_security" "testacc" {
		module_id                             = "%s"
		node_id                               = "Master"
		http_listener_context_path            = "/auth"
		http_listener_unhealthy_response_code = 503
		http_listener_port                    = 9999
		oidc_issuer_url                       = "http://keycloak:8080/auth/realms/lra"
		oidc_federate_mode_enabled            = true
		smart_authorization_enforce_approved_scopes = true
		smart_authorization_scopes_supported        = "launch fhirUser openid profile patient/*.read"
		sessions_max_concurrent_sessions_per_user   = 3
		dependency_fhir_persistence_module          = "PERSISTENCE_ALL"
}`, moduleName)
}
